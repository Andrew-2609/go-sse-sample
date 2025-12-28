import { useSSE } from '../hooks/useSSE'
import MetricCard from './MetricCard'
import ConnectionStatus from './ConnectionStatus'
import DebugPanel from './DebugPanel'
import './Dashboard.css'

function Dashboard() {
  const { metrics, connectionStatus, debugMessages, isHealthy } = useSSE()

  return (
    <div className="dashboard">
      <header className="dashboard-header">
        <h1>Metrics Dashboard</h1>
        <ConnectionStatus status={connectionStatus} />
      </header>

      {metrics.length === 0 ? (
        <div className="empty-state">
          <p>No metrics yet. Create a metric to see it appear here in real-time!</p>
          <p className="hint">Connect to the API and create metrics to see live updates via SSE.</p>
        </div>
      ) : (
        <div className="metrics-grid">
          {metrics.map((metric) => (
            <MetricCard key={metric.id} metric={metric} />
          ))}
        </div>
      )}

      <DebugPanel messages={debugMessages} isHealthy={isHealthy} />
    </div>
  )
}

export default Dashboard

