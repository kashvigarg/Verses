import { cookies } from "next/headers"
import { redirect } from "next/navigation"
import { MainLayout } from "@/components/main-layout"
import { UserProfile } from "@/components/user-profile"
import { getUserProfile } from "@/lib/api"

export default async function ProfilePage({ params }: { params: { username: string } }) {
  // const cookieStore = cookies()
  // const token = cookieStore.get("auth_token")

  // if (!token) {
  //   redirect("/login")
  // }

  const userData = {
    name: "John Doe",
    username: "johndoe123",
    id: "usr_001",
    follower: true,
    follows_back: false,
    followers: 250,
    following: 180,
    proses: [
        {
            id: "post_001",
            body: "Just started learning JavaScript!",
            created_at: "2025-03-23T08:30:00Z",
            updated_at: "2025-03-23T09:00:00Z",
            mine: true,
            liked: true,
            likes_count: 120,
            comments: 15
        },
        {
            id: "post_002",
            body: "Coding late at night hits different.",
            created_at: "2025-03-22T22:45:00Z",
            updated_at: "2025-03-23T00:00:00Z",
            mine: true,
            liked: false,
            likes_count: 75,
            comments: 8
        }
    ]
};

  // await getUserProfile(params.username, token.value)

  return (
    <MainLayout>
      <UserProfile user={userData} />
    </MainLayout>
  )
}

