import { useMemo, memo } from 'react'
import { LineChart, Line, XAxis, YAxis, CartesianGrid, Tooltip, ResponsiveContainer } from 'recharts'
import './MetricCard.css'

const MetricCard = memo(function MetricCard({ metric }) {
  // Memoize chart data - only recalculate when readings change
  const chartData = useMemo(() => {
    return metric.readings
      .slice(-20)
      .map((reading) => ({
        time: reading.timestamp.toLocaleTimeString(),
        value: reading.value,
        timestamp: reading.timestamp
      }))
  }, [metric.readings])

  // Memoize latest reading
  const latestReading = useMemo(() => {
    return metric.readings.length > 0 
      ? metric.readings[metric.readings.length - 1]
      : null
  }, [metric.readings])

  // Memoize average calculation
  const averageValue = useMemo(() => {
    if (metric.readings.length === 0) return 0
    return (metric.readings.reduce((sum, r) => sum + r.value, 0) / metric.readings.length).toFixed(2)
  }, [metric.readings])

  // Memoize recent readings list
  const recentReadings = useMemo(() => {
    return metric.readings.slice(-5).reverse()
  }, [metric.readings])

  return (
    <div className="metric-card">
      <div className="metric-card-header">
        <h2>{metric.name}</h2>
        <span className="metric-id">ID: {metric.id.slice(0, 8)}...</span>
      </div>

      <div className="metric-stats">
        <div className="stat-item">
          <span className="stat-label">Latest Value</span>
          <span className="stat-value">
            {latestReading ? latestReading.value.toFixed(2) : 'N/A'}
          </span>
        </div>
        <div className="stat-item">
          <span className="stat-label">Average</span>
          <span className="stat-value">{averageValue}</span>
        </div>
        <div className="stat-item">
          <span className="stat-label">Total Readings</span>
          <span className="stat-value">{metric.readings.length}</span>
        </div>
      </div>

      {metric.readings.length > 0 ? (
        <div className="metric-chart">
          <ResponsiveContainer width="100%" height={200}>
            <LineChart 
              data={chartData}
              margin={{ top: 5, right: 5, left: 5, bottom: 5 }}
            >
              <CartesianGrid strokeDasharray="3 3" stroke="#e0e0e0" />
              <XAxis 
                dataKey="time" 
                stroke="#666"
                fontSize={12}
                interval="preserveStartEnd"
              />
              <YAxis 
                stroke="#666"
                fontSize={12}
              />
              <Tooltip 
                contentStyle={{
                  backgroundColor: '#fff',
                  border: '1px solid #e0e0e0',
                  borderRadius: '4px'
                }}
                labelFormatter={(value) => `Time: ${value}`}
                formatter={(value) => [value.toFixed(2), 'Value']}
              />
              <Line 
                type="monotone" 
                dataKey="value" 
                stroke="#667eea" 
                strokeWidth={2}
                dot={false}
                isAnimationActive={true}
                animationDuration={300}
                activeDot={{ r: 4 }}
              />
            </LineChart>
          </ResponsiveContainer>
        </div>
      ) : (
        <div className="no-readings">
          <p>No readings yet. Add a reading to see the chart update in real-time!</p>
        </div>
      )}

      {metric.readings.length > 0 && (
        <div className="metric-readings-list">
          <h3>Recent Readings</h3>
          <div className="readings-scroll">
            {recentReadings.map((reading) => (
              <div key={reading.id} className="reading-item">
                <span className="reading-value">{reading.value.toFixed(2)}</span>
                <span className="reading-time">
                  {reading.timestamp.toLocaleString()}
                </span>
              </div>
            ))}
          </div>
        </div>
      )}
    </div>
  )
}, (prevProps, nextProps) => {
  // Return true if props are equal (skip re-render), false if different (re-render)
  
  // Re-render if ID or name changed
  if (prevProps.metric.id !== nextProps.metric.id || 
      prevProps.metric.name !== nextProps.metric.name) {
    return false // Props changed, should re-render
  }
  
  // Re-render if readings count changed
  if (prevProps.metric.readings.length !== nextProps.metric.readings.length) {
    return false // Props changed, should re-render
  }
  
  // Check if the last reading changed (most common case for updates)
  const prevLast = prevProps.metric.readings[prevProps.metric.readings.length - 1]
  const nextLast = nextProps.metric.readings[nextProps.metric.readings.length - 1]
  
  if (!prevLast && !nextLast) return true // Both empty, no change, skip re-render
  if (!prevLast || !nextLast) return false // One is empty, changed, should re-render
  if (prevLast.id !== nextLast.id || prevLast.value !== nextLast.value) {
    return false // Last reading changed, should re-render
  }
  
  // No significant changes detected, skip re-render
  return true
})

export default MetricCard

