export default {
  async fetch(request, env) {
    const url = new URL(request.url);

    if (url.pathname === "/api/notes") {
      const configuredBackendUrl =
        env && "BACKEND_URL" in env && typeof env.BACKEND_URL === "string"
          ? env.BACKEND_URL
          : undefined;

      if (!configuredBackendUrl) {
        console.error("BACKEND_URL not configured for worker", {
          requestHost: url.hostname,
        });
        return new Response(
          JSON.stringify({
            error: "backend_url_not_configured",
            message:
              "Set BACKEND_URL as a Worker var/secret. In dev, put it in frontend/.dev.vars.",
          }),
          { status: 500, headers: { "content-type": "application/json" } },
        );
      }

      try {
        const backendResponse = await fetch(
          `${configuredBackendUrl}/api/notes`,
          {
          headers: {
            accept: "application/json",
          },
          },
        );

        return new Response(backendResponse.body, {
          status: backendResponse.status,
          headers: {
            "content-type":
              backendResponse.headers.get("content-type") ?? "application/json",
          },
        });
      } catch (error) {
        console.error("Backend fetch failed", {
          configuredBackendUrl,
          error: error instanceof Error ? error.message : String(error),
        });
        return new Response(
          JSON.stringify({
            error: "backend_unreachable",
            backendBaseUrl: configuredBackendUrl,
          }),
          {
            status: 502,
            headers: { "content-type": "application/json" },
          },
        );
      }
    }

    if (url.pathname.startsWith("/api/")) {
      return Response.json({
        name: "Cloudflare",
      });
    }
		return new Response(null, { status: 404 });
  },
} satisfies ExportedHandler<Env>;
