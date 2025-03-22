// Server-side API functions

export async function getUserProfile(username: string, token: string) {
  const response = await fetch(`${process.env.NEXT_PUBLIC_API_URL || ""}/api/users/${username}`, {
    headers: {
      Authorization: `Bearer ${token}`,
    },
    next: { revalidate: 60 }, // Revalidate every 60 seconds
  })

  if (!response.ok) {
    throw new Error("Failed to fetch user profile")
  }

  return response.json()
}

export async function getProseById(proseId: string, token: string) {
  const response = await fetch(`${process.env.NEXT_PUBLIC_API_URL || ""}/api/prose/${proseId}`, {
    headers: {
      Authorization: `Bearer ${token}`,
    },
    next: { revalidate: 60 }, // Revalidate every 60 seconds
  })

  if (!response.ok) {
    throw new Error("Failed to fetch prose")
  }

  return response.json()
}

