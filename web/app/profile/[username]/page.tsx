import { cookies } from "next/headers"
import { redirect } from "next/navigation"
import { MainLayout } from "@/components/main-layout"
import { UserProfile } from "@/components/user-profile"
import { getUserProfile } from "@/lib/api"

export default async function ProfilePage({ params }: { params: { username: string } }) {
  const cookieStore = cookies()
  const token = cookieStore.get("auth_token")

  if (!token) {
    redirect("/login")
  }

  const userData = await getUserProfile(params.username, token.value)

  return (
    <MainLayout>
      <UserProfile user={userData} />
    </MainLayout>
  )
}

