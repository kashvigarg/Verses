import { redirect } from "next/navigation"
import { cookies } from "next/headers"
import { Timeline } from "@/components/timeline"
import { MainLayout } from "@/components/main-layout"

export default function Home() {
  // const cookieStore = cookies()
  // const token = cookieStore.get("auth_token")

  // if (!token) {
  //   redirect("/login")
  // }

  return (
    <MainLayout>
      <Timeline />
    </MainLayout>
  )
}

