// Auto-scroll log panel on new entries
const log = document.getElementById('log');
new MutationObserver(() => {
    if (!document.getElementById('log-panel').classList.contains('collapsed')) {
        log.scrollTop = log.scrollHeight;
    }
}).observe(log, { childList: true });

// Client-side car animation using time-based linear interpolation.
//
// Rather than lerping (which always lags behind), each car moves at constant
// velocity from its previous target to the current target over one estimated
// update interval, then extrapolates slightly past it to hide the gap before
// the next update arrives.
const SVG_NS = 'http://www.w3.org/2000/svg';

const targets = {};          // latest server position
const prevTargets = {};      // position at the moment the latest target arrived
const targetAt = {};         // performance.now() when latest target was received
const estInterval = {};      // per-car EMA of observed update interval (ms)
const current = {};          // current rendered position

function updateCarTargets(cars) {
    const now = performance.now();
    for (const [num, car] of Object.entries(cars)) {
        if (targets[num]) {
            const dt = now - (targetAt[num] || now);
            estInterval[num] = estInterval[num]
                ? estInterval[num] * 0.7 + dt * 0.3
                : dt;
        }
        prevTargets[num] = current[num] ? { ...current[num] } : car;
        targets[num] = car;
        targetAt[num] = now;
        if (!current[num]) current[num] = { x: car.x, y: car.y };
    }
}

function rafTick(now) {
    const carsG = document.getElementById('cars');
    if (carsG) {
        for (const [num, target] of Object.entries(targets)) {
            const prev = prevTargets[num];
            const cur = current[num];
            const elapsed = now - (targetAt[num] || now);
            const interval = estInterval[num] || 300;

            // alpha 0→1 = interpolating to target; slight overshoot fills the gap
            // before the next update. Starting from current (not old server target)
            // means overshoot unwinds gradually rather than slingshotting.
            const alpha = Math.min(elapsed / interval, 1.15);
            cur.x = prev.x + (target.x - prev.x) * alpha;
            cur.y = prev.y + (target.y - prev.y) * alpha;

            let circle = document.getElementById('car-' + num);
            if (!circle) {
                circle = document.createElementNS(SVG_NS, 'circle');
                circle.id = 'car-' + num;
                circle.setAttribute('r', '9');
                circle.setAttribute('stroke', '#111');
                circle.setAttribute('stroke-width', '1');
                const text = document.createElementNS(SVG_NS, 'text');
                text.id = 'car-label-' + num;
                text.setAttribute('text-anchor', 'middle');
                text.setAttribute('font-size', '7');
                text.setAttribute('font-family', 'monospace');
                text.setAttribute('font-weight', 'bold');
                carsG.appendChild(circle);
                carsG.appendChild(text);
            }

            const text = document.getElementById('car-label-' + num);
            const cx = cur.x.toFixed(1);
            const cy = cur.y.toFixed(1);
            circle.setAttribute('cx', cx);
            circle.setAttribute('cy', cy);
            circle.setAttribute('fill', target.off ? '#555555' : target.color);
            if (text) {
                text.setAttribute('x', cx);
                text.setAttribute('y', (cur.y + 2.5).toFixed(1));
                text.setAttribute('fill', target.off ? '#999999' : '#111111');
                text.textContent = target.label;
            }
        }
    }
    requestAnimationFrame(rafTick);
}

requestAnimationFrame(rafTick);