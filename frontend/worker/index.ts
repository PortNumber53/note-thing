export default {
  async fetch(request, env) {
    const url = new URL(request.url);

    // Proxy /api/*, /auth/*, and /callback/* to backend
    if (url.pathname.startsWith("/api/") || url.pathname.startsWith("/auth/") || url.pathname.startsWith("/callback/")) {
      const backendUrl =
        env && "BACKEND_URL" in env && typeof env.BACKEND_URL === "string"
          ? env.BACKEND_URL
          : undefined;

      if (!backendUrl) {
        return new Response(
          JSON.stringify({
            error: "backend_url_not_configured",
            message: "Set BACKEND_URL as a Worker var/secret. In dev, put it in frontend/.dev.vars.",
          }),
          { status: 500, headers: { "content-type": "application/json" } },
        );
      }

      const targetUrl = `${backendUrl}${url.pathname}${url.search}`;

      const headers = new Headers();
      headers.set("Accept", "application/json");

      // Forward relevant headers
      const authorization = request.headers.get("Authorization");
      if (authorization) headers.set("Authorization", authorization);

      const contentType = request.headers.get("Content-Type");
      if (contentType) headers.set("Content-Type", contentType);

      try {
        const backendResponse = await fetch(targetUrl, {
          method: request.method,
          headers,
          body: request.method !== "GET" && request.method !== "HEAD"
            ? request.body
            : undefined,
        });

        // Forward redirect responses as-is
        if (backendResponse.status >= 300 && backendResponse.status < 400) {
          return backendResponse;
        }

        return new Response(backendResponse.body, {
          status: backendResponse.status,
          headers: {
            "content-type": backendResponse.headers.get("content-type") ?? "application/json",
          },
        });
      } catch (error) {
        console.error("Backend fetch failed", {
          targetUrl,
          error: error instanceof Error ? error.message : String(error),
        });
        return new Response(
          JSON.stringify({ error: "backend_unreachable" }),
          { status: 502, headers: { "content-type": "application/json" } },
        );
      }
    }

    return new Response(null, { status: 404 });
  },
} satisfies ExportedHandler<Env>;
