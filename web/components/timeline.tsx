"use client"

import { useState, useEffect } from "react"
import { useAuth } from "@/lib/auth-hooks"
import { ProseCard } from "@/components/prose-card"
import { Skeleton } from "@/components/ui/skeleton"
import { useToast } from "@/components/ui/use-toast"
import { Alert, AlertDescription } from "@/components/ui/alert"
import { AlertCircle, RefreshCw } from "lucide-react"
import { Button } from "@/components/ui/button"
import { useSSE } from "@/lib/use-sse"

type TimelineItem = {
  id: number
  userid?: string
  prose: {
    id: string
    body: string
    created_at: string
    updated_at: string
    username: string
    mine: boolean
    liked: boolean
    likes_count: number
    comments: number
  }
}

export function Timeline() {
  const { user, token } = useAuth()
  const [timelineItems, setTimelineItems] = useState<TimelineItem[]>([])
  const [isLoading, setIsLoading] = useState(true)
  const [error, setError] = useState<string | null>(null)
  const { toast } = useToast()

  // Try to use SSE for timeline, fall back to regular fetch
  const { data, error: sseError } = useSSE<TimelineItem[]>("/api/timeline", token, {
    onMessage: (data) => {
      if (data) {
        setTimelineItems(data)
        setIsLoading(false)
      }
    },
    fallbackToFetch: true,
  })

  useEffect(() => {
    if (data) {
      setTimelineItems(data)
      setIsLoading(false)
    }

    if (sseError) {
      setError("Failed to load timeline. Please try again.")
      setIsLoading(false)
    }
  }, [data, sseError])

  const refreshTimeline = async () => {
    setIsLoading(true)
    setError(null)

    try {
      const response = await fetch("/api/timeline", {
        headers: {
          Authorization: `Bearer ${token}`,
        },
      })

      if (!response.ok) {
        throw new Error("Failed to refresh timeline")
      }

      const data = await response.json()
      setTimelineItems(data)
    } catch (err) {
      setError("Failed to refresh timeline. Please try again.")
      toast({
        title: "Error",
        description: "Failed to refresh timeline",
        variant: "destructive",
      })
    } finally {
      setIsLoading(false)
    }
  }

  const handleLikeToggle = async (proseId: string) => {
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

      // Update the prose in the list
      setTimelineItems(
        timelineItems.map((item) => {
          if (item.prose.id === proseId) {
            return {
              ...item,
              prose: {
                ...item.prose,
                liked: result.liked,
                likes_count: result.likes_count,
              },
            }
          }
          return item
        }),
      )
    } catch (err) {
      toast({
        title: "Error",
        description: "Failed to like/unlike verse",
        variant: "destructive",
      })
    }
  }

  if (isLoading) {
    return (
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
    )
  }

  if (error) {
    return (
      <Alert variant="destructive" className="mb-4">
        <AlertCircle className="h-4 w-4" />
        <AlertDescription>{error}</AlertDescription>
        <Button variant="outline" size="sm" className="ml-auto" onClick={refreshTimeline}>
          <RefreshCw className="mr-2 h-4 w-4" />
          Retry
        </Button>
      </Alert>
    )
  }

  return (
    <div className="space-y-6">
      <div className="flex items-center justify-between">
        <h2 className="text-2xl font-serif font-bold">Your Timeline</h2>
        <Button variant="outline" size="sm" onClick={refreshTimeline}>
          <RefreshCw className="mr-2 h-4 w-4" />
          Refresh
        </Button>
      </div>

      {timelineItems.length === 0 ? (
        <div className="rounded-lg border bg-white dark:bg-slate-900 p-8 text-center">
          <p className="text-muted-foreground">No verses found. Follow writers to see their work here.</p>
        </div>
      ) : (
        timelineItems.map((item) => (
          <ProseCard key={item.prose.id} prose={item.prose} onLikeToggle={() => handleLikeToggle(item.prose.id)} />
        ))
      )}
    </div>
  )
}

