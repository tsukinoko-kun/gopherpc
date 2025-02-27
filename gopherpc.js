globalThis.gopherpc = new Proxy(
  {},
  {
    get(_, property) {
      return async (...args) => {
        const resp = await fetch("/__gopherpc__/rpc", {
          method: "POST",
          headers: {
            Accept: "application/json",
            "Content-Type": "application/json",
          },
          body: JSON.stringify({ func_name: property, args }),
          cache: "no-cache",
        });
        const o = await resp.json();
        if ("type" in o) {
          switch (o.type) {
            case "error":
              throw new Error(o.error ?? o);
            case "ok":
              return o.result;
          }
        }
        throw new Error("unexpected response: " + JSON.stringify(o));
      };
    },
  },
);
