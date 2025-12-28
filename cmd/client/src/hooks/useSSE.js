import { useEffect, useState, useRef, useCallback, useMemo } from 'react'

const API_BASE_URL = 'http://localhost:8089'

export function useSSE() {
  const [metrics, setMetrics] = useState(new Map())
  const [connectionStatus, setConnectionStatus] = useState('disconnected')
  const [debugMessages, setDebugMessages] = useState([])
  const eventSourceRef = useRef(null)
  const updateQueueRef = useRef([])
  const rafIdRef = useRef(null)

  const addDebugMessage = useCallback((message, type = 'info') => {
    const timestamp = new Date().toLocaleTimeString()
    setDebugMessages((prev) => {
      const newMessages = [...prev, { message, type, timestamp }]
      // Keep only last 20 messages
      return newMessages.slice(-20)
    })
  }, [])

  // Batch state updates using requestAnimationFrame for smooth rendering
  const flushUpdates = useCallback(() => {
    if (rafIdRef.current) {
      cancelAnimationFrame(rafIdRef.current)
    }
    
    rafIdRef.current = requestAnimationFrame(() => {
      if (updateQueueRef.current.length === 0) return
      
      setMetrics((prev) => {
        const newMap = new Map(prev)
        for (const update of updateQueueRef.current) {
          update(newMap)
        }
        updateQueueRef.current = []
        return newMap
      })
    })
  }, [])

  useEffect(() => {
    let wasConnected = false
    let isCleaningUp = false
    let initialLoadComplete = false

    // Fetch initial state (all metrics with readings) before connecting to SSE
    const fetchInitialState = async () => {
      try {
        addDebugMessage('Fetching initial metrics state...', 'info')
        const response = await fetch(`${API_BASE_URL}/metrics?with_readings=true`)
        
        if (!response.ok) {
          throw new Error(`Failed to fetch initial state: ${response.statusText}`)
        }
        
        const metricsData = await response.json()
        
        // Transform the API response into our internal format
        const initialMetrics = new Map()
        for (const metric of metricsData) {
          const readings = (metric.readings || []).map(reading => ({
            id: reading.id,
            value: reading.value,
            timestamp: new Date(reading.timestamp)
          }))
          
          // Parse input_frequency from nanoseconds (Go time.Duration format)
          // Go sends time.Duration as nanoseconds in JSON
          const inputFrequencyMs = metric.input_frequency 
            ? (typeof metric.input_frequency === 'number' 
                ? metric.input_frequency / 1_000_000  // Convert nanoseconds to milliseconds
                : 0)
            : 0
          
          initialMetrics.set(metric.id, {
            id: metric.id,
            name: metric.name,
            inputFrequencyMs: inputFrequencyMs,
            readings: readings.sort((a, b) => a.timestamp - b.timestamp)
          })
        }
        
        setMetrics(initialMetrics)
        initialLoadComplete = true
        addDebugMessage(`Loaded ${initialMetrics.size} metric(s) with existing readings`, 'success')
      } catch (error) {
        addDebugMessage(`Error fetching initial state: ${error.message}`, 'error')
        // Continue anyway - SSE will populate as events come in
        initialLoadComplete = true
      }
    }

    // Clean up any existing connection first
    if (eventSourceRef.current) {
      eventSourceRef.current.close()
    }

    // Fetch initial state, then set up SSE connection
    fetchInitialState().then(() => {
      if (isCleaningUp) return
      
      // Set up SSE connection after initial state is loaded
      const eventSource = new EventSource(`${API_BASE_URL}/events/watch`)
      eventSourceRef.current = eventSource

      // Handle connection open
      eventSource.onopen = () => {
        if (isCleaningUp) return
        wasConnected = true
        setConnectionStatus('connected')
        addDebugMessage('SSE connection opened', 'success')
      }

      // Handle errors - EventSource fires onerror even on successful connections
      // Only set error if connection is actually closed AND we were previously connected
      eventSource.onerror = () => {
        if (isCleaningUp) return
        
        const readyState = eventSource.readyState
        
        // 0 = CONNECTING, 1 = OPEN, 2 = CLOSED
        if (readyState === 2 && wasConnected) {
          // Connection was open but is now closed
          setConnectionStatus('error')
          addDebugMessage('SSE connection closed unexpectedly', 'error')
          wasConnected = false
        } else if (readyState === 1) {
          // Connection is open - transient error, keep as connected
          setConnectionStatus('connected')
        }
        // If readyState === 0 (CONNECTING), don't change status
      }

      // Handle metric_created events
      eventSource.addEventListener('metric_created', (event) => {
        if (isCleaningUp) return
        
        try {
          const data = JSON.parse(event.data)
          addDebugMessage(`Metric created: ${data.name} (${data.id.slice(0, 8)}...)`, 'success')
          setMetrics((prev) => {
            const newMap = new Map(prev)
            // Only add if it doesn't already exist (avoid duplicates from initial load)
            if (!newMap.has(data.id)) {
              // Parse input_frequency from nanoseconds (Go time.Duration format)
              const inputFrequencyMs = data.input_frequency 
                ? (typeof data.input_frequency === 'number' 
                    ? data.input_frequency / 1_000_000  // Convert nanoseconds to milliseconds
                    : 0)
                : 0
              
              newMap.set(data.id, {
                id: data.id,
                name: data.name,
                inputFrequencyMs: inputFrequencyMs,
                readings: []
              })
            }
            return newMap
          })
        } catch (error) {
          addDebugMessage(`Error parsing metric_created event: ${error.message}`, 'error')
        }
      })

      // Handle metric_reading_created events - optimized for high frequency
      eventSource.addEventListener('metric_reading_created', (event) => {
        if (isCleaningUp) return
        
        try {
          const data = JSON.parse(event.data)
          
          // Queue the update instead of applying immediately
          updateQueueRef.current.push((newMap) => {
            const metric = newMap.get(data.metric_id)
            
            if (metric) {
              const newReading = {
                id: data.id,
                value: data.value,
                timestamp: new Date(data.timestamp)
              }
              
              // Check if reading already exists (avoid duplicates)
              const readingExists = metric.readings.some(r => r.id === data.id)
              if (!readingExists) {
                // Optimize: append to end if timestamp is newer (most common case)
                const lastReading = metric.readings[metric.readings.length - 1]
                const isNewer = !lastReading || newReading.timestamp >= lastReading.timestamp
                
                const updatedMetric = {
                  ...metric,
                  readings: isNewer
                    ? [...metric.readings, newReading] // Fast path: just append
                    : [...metric.readings, newReading].sort((a, b) => a.timestamp - b.timestamp) // Slow path: sort
                }
                
                newMap.set(data.metric_id, updatedMetric)
              }
            }
          })
          
          // Flush updates on next animation frame for smooth rendering
          flushUpdates()
        } catch (error) {
          addDebugMessage(`Error parsing metric_reading_created event: ${error.message}`, 'error')
        }
      })

      // Handle connection messages (generic messages)
      eventSource.onmessage = (event) => {
        if (isCleaningUp) return
        
        if (event.data === 'connected') {
          wasConnected = true
          setConnectionStatus('connected')
          addDebugMessage('SSE connection confirmed', 'success')
        } else if (event.data === 'disconnected') {
          wasConnected = false
          setConnectionStatus('disconnected')
          addDebugMessage('SSE disconnected', 'warning')
        }
      }
    })

    return () => {
      isCleaningUp = true
      if (rafIdRef.current) {
        cancelAnimationFrame(rafIdRef.current)
      }
      if (eventSourceRef.current) {
        eventSourceRef.current.close()
        setConnectionStatus('disconnected')
      }
    }
  }, [addDebugMessage, flushUpdates])

  const isHealthy = connectionStatus === 'connected' && debugMessages.filter(m => m.type === 'error').length === 0

  // Memoize metrics array - convert Map to array efficiently
  // Since we batch updates with requestAnimationFrame, the Map reference only changes on flush
  const metricsArray = useMemo(() => {
    return Array.from(metrics.values())
  }, [metrics])

  return {
    metrics: metricsArray,
    connectionStatus,
    debugMessages,
    isHealthy
  }
}

