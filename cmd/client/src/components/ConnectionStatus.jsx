import './ConnectionStatus.css'

function ConnectionStatus({ status }) {
  const getStatusClass = () => {
    switch (status) {
      case 'connected':
        return 'status-connected'
      case 'error':
        return 'status-error'
      default:
        return 'status-disconnected'
    }
  }

  const getStatusText = () => {
    switch (status) {
      case 'connected':
        return 'Connected'
      case 'error':
        return 'Connection Error'
      default:
        return 'Disconnected'
    }
  }

  return (
    <div className={`connection-status ${getStatusClass()}`}>
      <span className="status-dot"></span>
      <span className="status-text">{getStatusText()}</span>
    </div>
  )
}

export default ConnectionStatus

