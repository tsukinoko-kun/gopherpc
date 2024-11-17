const url = new URL(window.location.href);
url.protocol = url.protocol === 'https:' ? 'wss:' : 'ws:';
url.pathname = url.pathname.split('__gopherpc__')[0] + '__gopherpc__/ws';

const ws = new WebSocket(url.href);
const respQueue = new Map();

globalThis.gopherpc = new Proxy({}, {
    get(_, property) {
        return (...args) => {
            const id = gopherpcCallId();
            const expProm = explodedPromise();
            respQueue.set(id, expProm);
            ws.send(JSON.stringify({ func_name: property, args, id }));
            return expProm.promise;
        };
    }
});

function gopherpcCallId() {
    return Math.random().toString(36).substring(2);
}

function explodedPromise() {
    let resolve, reject;
    const promise = new Promise((res, rej) => {
        resolve = res;
        reject = rej;
    });
    return { promise, resolve, reject };
}

ws.onopen = () => {
    console.debug('GopheRPC connected');
};

ws.onmessage = (event) => {
    const { id, result, error } = JSON.parse(event.data);
    const { resolve, reject } = respQueue.get(id);
    if (error) {
        reject(new Error(error));
    } else {
        resolve(result);
    }
    respQueue.delete(id);
};

ws.onclose = () => {
    console.debug('GopheRPC disconnected');
};
