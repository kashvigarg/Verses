import { cookies } from "next/headers"
import { redirect } from "next/navigation"
import { MainLayout } from "@/components/main-layout"
import { ProseDetail } from "@/components/prose-detail"
import { getProseById } from "@/lib/api"

export default async function ProsePage({ params }: { params: { proseId: string } }) {
  // const cookieStore = cookies()
  // const token = cookieStore.get("auth_token")

  // if (!token) {
  //   redirect("/login")
  // }

  const proseData = {
    id: "post_001",
    body: "Exploring the world of coding!",
    created_at: "2025-03-23T12:00:00Z",
    updated_at: "2025-03-23T12:30:00Z",
    username: "johndoe123",
    mine: true,
    liked: true,
    likes_count: 42,
    comments: 10
};
  // await getProseById(params.proseId, token.value)

  return (
    <MainLayout>
      <ProseDetail prose={proseData} />
    </MainLayout>
  )
}

