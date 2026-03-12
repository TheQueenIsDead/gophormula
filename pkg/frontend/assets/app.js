    const log = document.getElementById('log');
    new MutationObserver(() => {
    if (!document.getElementById('log-panel').classList.contains('collapsed')) {
    log.scrollTop = log.scrollHeight;
}
}).observe(log, { childList: true });
