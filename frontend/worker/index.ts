export default {
  async fetch(request, env) {
    const url = new URL(request.url);

    if (url.pathname === "/api/notes") {
      const backendBaseUrl =
        (env && "BACKEND_URL" in env && env.BACKEND_URL) || "http://localhost:18611";

      const backendResponse = await fetch(`${backendBaseUrl}/api/notes`, {
        headers: {
          accept: "application/json",
        },
      });

      return new Response(backendResponse.body, {
        status: backendResponse.status,
        headers: {
          "content-type": backendResponse.headers.get("content-type") ?? "application/json",
        },
      });
    }

    if (url.pathname.startsWith("/api/")) {
      return Response.json({
        name: "Cloudflare",
      });
    }
		return new Response(null, { status: 404 });
  },
} satisfies ExportedHandler<Env>;
