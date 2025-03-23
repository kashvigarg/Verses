import { cookies } from "next/headers"
import { redirect } from "next/navigation"
import { MainLayout } from "@/components/main-layout"
import { NotificationsList } from "@/components/notifications-list"

export default function NotificationsPage() {
  // const cookieStore = cookies()
  // const token = cookieStore.get("auth_token")

  // if (!token) {
  //   redirect("/login")
  // }

  return (
    <MainLayout>
      <NotificationsList />
    </MainLayout>
  )
}

