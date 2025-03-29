import { cookies } from "next/headers"
import { redirect } from "next/navigation"
import { MainLayout } from "@/components/main-layout"
import { NotificationsList } from "@/components/notifications-list"

export default async function NotificationsPage() {
  const cookieStore = await cookies()
  const token = cookieStore.get("auth_token")

  if (!token) {
    redirect("/login")
  }

  return (
    <MainLayout>
      <NotificationsList />
    </MainLayout>
  )
}

