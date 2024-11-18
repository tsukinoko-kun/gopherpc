const respQueue = new Map();

function explodedPromise() {
    let resolve, reject;
    const promise = new Promise((res, rej) => {
        resolve = res;
        reject = rej;
    });
    return { promise, resolve, reject };
}

function newWs() {
    const url = new URL(window.location.href);
    url.protocol = url.protocol === 'https:' ? 'wss:' : 'ws:';
    url.pathname = '/__gopherpc__/ws';

    console.debug('GopheRPC connecting', url.href);

    const prom = explodedPromise();
    const ws = new WebSocket(url.href);

    ws.addEventListener("open", () => {
        prom.resolve();
        console.debug('GopheRPC connected');
    });

    ws.addEventListener("close", () => {
        console.debug('GopheRPC disconnected');
    });

    ws.addEventListener("error", (event) => {
        prom.reject(event);
        console.error('GopheRPC error', event);
    });

    ws.addEventListener("message", (event) => {
        const { id, result, error } = JSON.parse(event.data);
        const { resolve, reject } = respQueue.get(id);
        if (error) {
            reject(new Error(error));
        } else {
            resolve(result);
        }
        respQueue.delete(id);
    });

    if (ws.readyState === WebSocket.OPEN) {
        prom.resolve();
    }

    return { ws, prom: prom.promise };
}

function gopherpcCallId() {
    return Math.random().toString(36).substring(2);
}

let { ws, prom } = newWs();

window.addEventListener("beforeunload", () => {
    ws.close();
});

globalThis.goWs = () => ws;
globalThis.gopherpc = new Proxy({}, {
    get(_, property) {
        return async (...args) => {
            switch (ws.readyState) {
                case WebSocket.CLOSING:
                case WebSocket.CLOSED:
                    ({ ws, prom } = newWs());
                case WebSocket.CONNECTING:
                    await prom;
                    break;
            }
            const id = gopherpcCallId();
            const expProm = explodedPromise();
            respQueue.set(id, expProm);
            ws.send(JSON.stringify({ func_name: property, args, id }));
            return await expProm.promise;
        };
    }
});
