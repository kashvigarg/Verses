import { type NextRequest, NextResponse } from "next/server"

// This is a proxy API route that forwards requests to the backend
export async function GET(request: NextRequest, { params }: { params: { path: string[] } }) {
  // const path = params.path.join("/")
  const url = new URL(request.url);
  const path = url.pathname.replace("/api/", ""); 
  const { searchParams } = new URL(request.url)
  const token = searchParams.get("token") || request.headers.get("Authorization")?.split(" ")[1]

  // Check if this is an SSE endpoint
  const isSSE = path === "timeline" || path === "notifications" || path.includes("/comments")

  if (isSSE) {
    try {
      // Try SSE first
      const response = await fetch(`${process.env.API_URL}/api/${path}`, {
        headers: {
          Authorization: `Bearer ${token}`,
          Accept: "text/event-stream",
        },
      })

      if (response.ok) {
        // Set up SSE response
        const encoder = new TextEncoder()
        const stream = new ReadableStream({
          async start(controller) {
            const data = await response.json()

            // Send the initial data
            controller.enqueue(encoder.encode(`data: ${JSON.stringify(data)}\n\n`))

            // Keep the connection open
            const interval = setInterval(() => {
              controller.enqueue(encoder.encode(": keepalive\n\n"))
            }, 30000)

            // Clean up on close
            request.signal.addEventListener("abort", () => {
              clearInterval(interval)
              controller.close()
            })
          },
        })

        return new NextResponse(stream, {
          headers: {
            "Content-Type": "text/event-stream",
            "Cache-Control": "no-cache",
            Connection: "keep-alive",
          },
        })
      } else {
        // Fall back to regular response
        return NextResponse.json(await response.json(), { status: response.status })
      }
    } catch (error) {
      console.error("SSE error:", error)
      return NextResponse.json({ error: "Failed to connect to SSE" }, { status: 500 })
    }
  }

  // Regular API request
  try {
    const response = await fetch(`${process.env.API_URL}/api/${path}`, {
      headers: {
        Authorization: token ? `Bearer ${token}` : "",
        "Content-Type": request.headers.get("Content-Type") || "application/json",
      },
    })

    const data = await response.json()
    return NextResponse.json(data, { status: response.status })
  } catch (error) {
    console.error("API error:", error)
    return NextResponse.json({ error: "Failed to fetch data" }, { status: 500 })
  }
}

export async function POST(request: NextRequest, { params }: { params: { path: string[] } }) {

  const url = new URL(request.url);
  const path = url.pathname.replace("/api/", ""); // Remove "/api/" prefix

  // if (!path) {
  //   return new Response("Invalid request", { status: 400 });
  // }
  //const path = params.path.join("/")

  const token = request.headers.get("Authorization")?.split(" ")[1]
  const body = await request.json()

  try {
    const response = await fetch(`${process.env.API_URL}/api/${path}`, {
      method: "POST",
      headers: {
        Authorization: token ? `Bearer ${token}` : "",
        "Content-Type": "application/json",
      },
      body: JSON.stringify(body),
    })

    const data = await response.json()
    return NextResponse.json(data, { status: response.status })
  } catch (error) {
    console.error("API error:", error)
    return NextResponse.json({ error: "Failed to post data" }, { status: 500 })
  }
}

export async function DELETE(request: NextRequest, { params }: { params: { path: string[] } }) {
  const url = new URL(request.url);
  const path = url.pathname.replace("/api/", ""); 
  // const path = params.path.join("/")
  const token = request.headers.get("Authorization")?.split(" ")[1]

  try {
    const response = await fetch(`${process.env.API_URL}/api/${path}`, {
      method: "DELETE",
      headers: {
        Authorization: token ? `Bearer ${token}` : "",
      },
    })

    if (response.status === 204) {
      return new NextResponse(null, { status: 204 })
    }

    const data = await response.json()
    return NextResponse.json(data, { status: response.status })
  } catch (error) {
    console.error("API error:", error)
    return NextResponse.json({ error: "Failed to delete data" }, { status: 500 })
  }
}

export async function PUT(request: NextRequest, { params }: { params: { path: string[] } }) {
  // const path = params.path.join("/")
  const url = new URL(request.url);
  const path = url.pathname.replace("/api/", ""); 
  const token = request.headers.get("Authorization")?.split(" ")[1]
  const body = await request.json()

  try {
    const response = await fetch(`${process.env.API_URL}/api/${path}`, {
      method: "PUT",
      headers: {
        Authorization: token ? `Bearer ${token}` : "",
        "Content-Type": "application/json",
      },
      body: JSON.stringify(body),
    })

    const data = await response.json()
    return NextResponse.json(data, { status: response.status })
  } catch (error) {
    console.error("API error:", error)
    return NextResponse.json({ error: "Failed to update data" }, { status: 500 })
  }
}

