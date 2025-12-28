import './DebugPanel.css'

function DebugPanel({ messages, isHealthy }) {
  // Hide panel when healthy and no errors/warnings
  const hasErrorsOrWarnings = messages.some(msg => msg.type === 'error' || msg.type === 'warning')
  
  if (isHealthy && !hasErrorsOrWarnings) {
    return null
  }

  return (
    <div className={`debug-panel ${isHealthy ? 'debug-panel-healthy' : 'debug-panel-error'}`}>
      <div className="debug-panel-header">
        <h3>Debug Information</h3>
        <span className="debug-status">{isHealthy ? '✓ Healthy' : '⚠ Issues Detected'}</span>
      </div>
      <div className="debug-messages">
        {messages.length === 0 ? (
          <p className="debug-empty">No issues detected</p>
        ) : (
          messages.map((msg, index) => (
            <div key={index} className={`debug-message ${msg.type}`}>
              <span className="debug-timestamp">{msg.timestamp}</span>
              <span className="debug-content">{msg.message}</span>
            </div>
          ))
        )}
      </div>
    </div>
  )
}

export default DebugPanel

