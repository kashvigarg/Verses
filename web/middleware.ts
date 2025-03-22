import { NextResponse } from "next/server"
import type { NextRequest } from "next/server"

export function middleware(request: NextRequest) {
  const path = request.nextUrl.pathname

  // Define public paths that don't require authentication
  const isPublicPath = path === "/login" || path === "/signup"

  // Get the token from the cookies
  const token = request.cookies.get("auth_token")?.value

  // Redirect to login if accessing a protected route without a token
  if (!isPublicPath && !token) {
    return NextResponse.redirect(new URL("/login", request.url))
  }

  // Redirect to home if accessing login/signup with a valid token
  if (isPublicPath && token) {
    return NextResponse.redirect(new URL("/", request.url))
  }

  return NextResponse.next()
}

// Configure the paths that should trigger this middleware
export const config = {
  matcher: [
    /*
     * Match all request paths except:
     * - api routes
     * - static files (images, js, css, etc.)
     * - favicon.ico
     */
    "/((?!api|_next/static|_next/image|favicon.ico).*)",
  ],
}

