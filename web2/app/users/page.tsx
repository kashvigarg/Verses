import { cookies } from "next/headers"
import { redirect } from "next/navigation"
import { MainLayout } from "@/components/main-layout"
import { UsersList } from "@/components/users-list"

export default async function UsersPage() {
  const cookieStore = await cookies()
  const token = cookieStore.get("auth_token")

  if (!token) {
    redirect("/login")
  }

  return (
    <MainLayout>
      <UsersList />
    </MainLayout>
  )
}

