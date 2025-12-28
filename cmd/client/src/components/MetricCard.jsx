import { LineChart, Line, XAxis, YAxis, CartesianGrid, Tooltip, ResponsiveContainer } from 'recharts'
import './MetricCard.css'

function MetricCard({ metric }) {
  // Prepare data for the chart (last 20 readings for performance)
  const chartData = metric.readings
    .slice(-20)
    .map((reading) => ({
      time: reading.timestamp.toLocaleTimeString(),
      value: reading.value,
      timestamp: reading.timestamp
    }))

  const latestReading = metric.readings.length > 0 
    ? metric.readings[metric.readings.length - 1]
    : null

  const averageValue = metric.readings.length > 0
    ? (metric.readings.reduce((sum, r) => sum + r.value, 0) / metric.readings.length).toFixed(2)
    : 0

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
            <LineChart data={chartData}>
              <CartesianGrid strokeDasharray="3 3" stroke="#e0e0e0" />
              <XAxis 
                dataKey="time" 
                stroke="#666"
                fontSize={12}
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
                dot={{ fill: '#667eea', r: 3 }}
                activeDot={{ r: 5 }}
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
            {metric.readings.slice(-5).reverse().map((reading) => (
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
}

export default MetricCard

