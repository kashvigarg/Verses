"use client"

import { useState } from "react"
import { useRouter } from "next/navigation"
import { useAuth } from "@/lib/auth-hooks"
import { Avatar, AvatarFallback, AvatarImage } from "@/components/ui/avatar"
import { Button } from "@/components/ui/button"
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from "@/components/ui/card"
import { Tabs, TabsContent, TabsList, TabsTrigger } from "@/components/ui/tabs"
import { ProseCard } from "@/components/prose-card"
import { Skeleton } from "@/components/ui/skeleton"
import { useToast } from "@/components/ui/use-toast"
import { Loader2 } from "lucide-react"

type UserProfileProps = {
  user: {
    name: string
    username: string
    id: string
    follower: boolean
    follows_back: boolean
    followers: number
    following: number
    proses?: Array<{
      id: string
      body: string
      created_at: string
      updated_at: string
      mine: boolean
      liked: boolean
      likes_count: number
      comments: number
    }>
  }
}

const dummy_user = {
  name: "John Doe",
  username: "johndoe123",
  id: "usr_001",
  follower: true,
  follows_back: false,
  followers: 150,
  following: 100,
  proses: [
      {
          id: "post_001",
          body: "This is my first post!",
          created_at: "2025-03-23T10:00:00Z",
          updated_at: "2025-03-23T10:00:00Z",
          mine: true,
          liked: false,
          likes_count: 5,
          comments: 2
      },
      {
          id: "post_002",
          body: "Loving the JavaScript vibes!",
          created_at: "2025-03-22T15:30:00Z",
          updated_at: "2025-03-22T16:00:00Z",
          mine: true,
          liked: true,
          likes_count: 20,
          comments: 5
      }
  ]
};

const token = "dummy token";


export function UserProfile({ user: initialUser }: UserProfileProps) {
  const [user, setUser] = useState(initialUser)
  const [isLoading, setIsLoading] = useState(false)
  const { user: currentUser, token } = useAuth()
  const { toast } = useToast()
  const router = useRouter()

  const isCurrentUser = currentUser?.username === user.username

  const handleToggleFollow = async () => {
    setIsLoading(true)

    try {
      const response = await fetch(`/api/users/${user.username}/toggle_follow`, {
        method: "POST",
        headers: {
          Authorization: `Bearer ${token}`,
        },
      })

      if (!response.ok) {
        throw new Error("Failed to toggle follow")
      }

      const result = await response.json()

      setUser({
        ...user,
        follower: result.followed,
        followers: result.followers_count,
      })
    } catch (err) {
      toast({
        title: "Error",
        description: "Failed to follow/unfollow user",
        variant: "destructive",
      })
    } finally {
      setIsLoading(false)
    }
  }

  const handleLikeToggle = async (proseId: string) => {
    if (!user.proses) return

    try {
      const response = await fetch(`/api/prose/${proseId}/togglelike`, {
        method: "POST",
        headers: {
          Authorization: `Bearer ${token}`,
        },
      })

      if (!response.ok) {
        throw new Error("Failed to toggle like")
      }

      const result = await response.json()

      // Update the prose in the user's proses list
      setUser({
        ...user,
        proses: user.proses.map((prose) => {
          if (prose.id === proseId) {
            return {
              ...prose,
              liked: result.liked,
              likes_count: result.likes_count,
            }
          }
          return prose
        }),
      })
    } catch (err) {
      toast({
        title: "Error",
        description: "Failed to like/unlike verse",
        variant: "destructive",
      })
    }
  }

  const handleDeleteProse = async (proseId: string) => {
    if (!user.proses) return

    try {
      const response = await fetch(`/api/prose/${proseId}`, {
        method: "DELETE",
        headers: {
          Authorization: `Bearer ${token}`,
        },
      })

      if (!response.ok) {
        throw new Error("Failed to delete verse")
      }

      // Remove the prose from the user's proses list
      setUser({
        ...user,
        proses: user.proses.filter((prose) => prose.id !== proseId),
      })

      toast({
        title: "Verse deleted",
        description: "Your verse has been deleted successfully",
      })
    } catch (err) {
      toast({
        title: "Error",
        description: "Failed to delete verse",
        variant: "destructive",
      })
    }
  }

  return (
    <div className="space-y-6">
      <Card className="bg-white dark:bg-slate-900 border shadow-sm">
        <CardHeader className="pb-2">
          <div className="flex items-center justify-between">
            <div className="flex items-center gap-4">
              <Avatar className="h-16 w-16">
                <AvatarImage src={`https://avatar.vercel.sh/${user.username}`} />
                <AvatarFallback>{user.username.slice(0, 2).toUpperCase()}</AvatarFallback>
              </Avatar>

              <div>
                <CardTitle className="text-xl font-serif">{user.name}</CardTitle>
                <CardDescription>@{user.username}</CardDescription>
              </div>
            </div>

            {!isCurrentUser && (
              <Button variant={user.follower ? "outline" : "default"} onClick={handleToggleFollow} disabled={isLoading}>
                {isLoading ? <Loader2 className="h-4 w-4 animate-spin" /> : user.follower ? "Unfollow" : "Follow"}
              </Button>
            )}
          </div>
        </CardHeader>

        <CardContent className="space-y-4">
          <div className="flex gap-4 text-sm">
            <div>
              <span className="font-semibold">{user.followers}</span>{" "}
              <span className="text-muted-foreground">Followers</span>
            </div>
            <div>
              <span className="font-semibold">{user.following}</span>{" "}
              <span className="text-muted-foreground">Following</span>
            </div>
          </div>
        </CardContent>
      </Card>

      <Tabs defaultValue="verses">
        <TabsList className="w-full bg-white dark:bg-slate-900">
          <TabsTrigger value="verses" className="flex-1">
            Verses
          </TabsTrigger>
          <TabsTrigger value="likes" className="flex-1">
            Likes
          </TabsTrigger>
        </TabsList>

        <TabsContent value="verses" className="mt-4 space-y-6">
          {!user.proses ? (
            <div className="space-y-6">
              {[1, 2, 3].map((i) => (
                <div key={i} className="rounded-lg border bg-white dark:bg-slate-900 p-4 space-y-4">
                  <div className="flex items-center space-x-4">
                    <Skeleton className="h-12 w-12 rounded-full" />
                    <div className="space-y-2">
                      <Skeleton className="h-4 w-[150px]" />
                      <Skeleton className="h-4 w-[100px]" />
                    </div>
                  </div>
                  <Skeleton className="h-24 w-full" />
                  <div className="flex space-x-4">
                    <Skeleton className="h-8 w-16" />
                    <Skeleton className="h-8 w-16" />
                  </div>
                </div>
              ))}
            </div>
          ) : user.proses.length === 0 ? (
            <div className="rounded-lg border bg-white dark:bg-slate-900 p-8 text-center">
              <p className="text-muted-foreground">No verses found.</p>
            </div>
          ) : (
            user.proses.map((prose) => (
              <ProseCard
                key={prose.id}
                prose={{
                  ...prose,
                  username: user.username,
                }}
                onLikeToggle={() => handleLikeToggle(prose.id)}
                onDelete={isCurrentUser ? () => handleDeleteProse(prose.id) : undefined}
              />
            ))
          )}
        </TabsContent>

        <TabsContent value="likes" className="mt-4">
          <div className="rounded-lg border bg-white dark:bg-slate-900 p-8 text-center">
            <p className="text-muted-foreground">Liked verses will appear here.</p>
          </div>
        </TabsContent>
      </Tabs>
    </div>
  )
}

