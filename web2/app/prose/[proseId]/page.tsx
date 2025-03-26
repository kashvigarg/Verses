import { cookies } from "next/headers"
import { redirect } from "next/navigation"
import { MainLayout } from "@/components/main-layout"
import { ProseDetail } from "@/components/prose-detail"
import { getProseById } from "@/lib/api"

export default async function ProsePage({ params }: { params: { proseId: string } }) {
  const cookieStore = cookies()
  const token = cookieStore.get("auth_token")

  if (!token) {
    redirect("/login")
  }

  const proseData = await getProseById(params.proseId, token.value)

  return (
    <MainLayout>
      <ProseDetail prose={proseData} />
    </MainLayout>
  )
}

