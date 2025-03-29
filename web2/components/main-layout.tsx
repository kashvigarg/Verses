"use client"

import type React from "react"

import { useState } from "react"
import Link from "next/link"
import { usePathname } from "next/navigation"
import { useAuth } from "@/lib/auth-hooks"
import { Button } from "@/components/ui/button"
import { Home, Bell, User, LogOut, Menu, X, PenSquare, BookOpen, Users } from "lucide-react"
import { Sheet, SheetContent, SheetTrigger } from "@/components/ui/sheet"
import { NotificationIndicator } from "@/components/notification-indicator"
import { ComposeProseDialog } from "@/components/compose-prose-dialog"

export function MainLayout({ children }: { children: React.ReactNode }) {
  const { user, logout } = useAuth()
  const pathname = usePathname()
  const [isComposeOpen, setIsComposeOpen] = useState(false)
  const [isMobileMenuOpen, setIsMobileMenuOpen] = useState(false)

  const navItems = [
    { href: "/", label: "Home", icon: Home },
    { href: "/notifications", label: "Notifications", icon: Bell, indicator: true },
    { href: `/profile/${user?.username}`, label: "Profile", icon: User },
    {href: "/users", label: "Users", icon: Users}
  ]

  return (
    <div className="flex min-h-screen bg-[#f8f5f0] dark:bg-slate-950">
      {/* Desktop Sidebar */}
      <div className="hidden w-64 border-r bg-white dark:bg-slate-900 p-4 md:flex md:flex-col">
        <div className="mb-8 flex items-center">
          <BookOpen className="h-6 w-6 text-primary mr-2" />
          <h1 className="text-2xl font-serif font-bold text-primary">Verses</h1>
        </div>

        <nav className="flex flex-1 flex-col gap-2">
          {navItems.map((item) => (
            <Link key={item.href} href={item.href}>
              <Button variant={pathname === item.href ? "default" : "ghost"} className="w-full justify-start">
                <item.icon className="mr-2 h-5 w-5" />
                {item.label}
                {item.indicator && <NotificationIndicator className="ml-2" />}
              </Button>
            </Link>
          ))}

          <Button onClick={() => setIsComposeOpen(true)} className="mt-4">
            <PenSquare className="mr-2 h-5 w-5" />
            Write Verse
          </Button>
        </nav>

        <div className="mt-auto pt-4">
          <Button variant="outline" className="w-full justify-start" onClick={logout}>
            <LogOut className="mr-2 h-5 w-5" />
            Sign out
          </Button>
        </div>
      </div>

      {/* Mobile Header */}
      <div className="flex w-full flex-col md:ml-64 md:w-auto">
        <header className="sticky top-0 z-10 flex h-14 items-center border-b bg-white dark:bg-slate-900 px-4 md:hidden">
          <Sheet open={isMobileMenuOpen} onOpenChange={setIsMobileMenuOpen}>
            <SheetTrigger asChild>
              <Button variant="ghost" size="icon">
                <Menu className="h-5 w-5" />
              </Button>
            </SheetTrigger>
            <SheetContent side="left" className="w-64 p-0">
              <div className="flex h-14 items-center border-b px-4">
                <BookOpen className="h-5 w-5 text-primary mr-2" />
                <h1 className="text-xl font-serif font-bold text-primary">Verses</h1>
                <Button variant="ghost" size="icon" className="ml-auto" onClick={() => setIsMobileMenuOpen(false)}>
                  <X className="h-5 w-5" />
                </Button>
              </div>

              <nav className="flex flex-col gap-2 p-4">
                {navItems.map((item) => (
                  <Link key={item.href} href={item.href} onClick={() => setIsMobileMenuOpen(false)}>
                    <Button variant={pathname === item.href ? "default" : "ghost"} className="w-full justify-start">
                      <item.icon className="mr-2 h-5 w-5" />
                      {item.label}
                      {item.indicator && <NotificationIndicator className="ml-2" />}
                    </Button>
                  </Link>
                ))}

                <Button
                  onClick={() => {
                    setIsComposeOpen(true)
                    setIsMobileMenuOpen(false)
                  }}
                  className="mt-4"
                >
                  <PenSquare className="mr-2 h-5 w-5" />
                  Write Verse
                </Button>

                <div className="mt-auto pt-4">
                  <Button variant="outline" className="w-full justify-start" onClick={logout}>
                    <LogOut className="mr-2 h-5 w-5" />
                    Sign out
                  </Button>
                </div>
              </nav>
            </SheetContent>
          </Sheet>

          <div className="flex items-center ml-4">
            <BookOpen className="h-5 w-5 text-primary mr-2" />
            <h1 className="text-xl font-serif font-bold text-primary">Verses</h1>
          </div>

          <Button variant="ghost" size="icon" className="ml-auto" onClick={() => setIsComposeOpen(true)}>
            <PenSquare className="h-5 w-5" />
          </Button>
        </header>

        <main className="flex-1 p-4 max-w-3xl mx-auto w-full">{children}</main>
      </div>

      <ComposeProseDialog open={isComposeOpen} onOpenChange={setIsComposeOpen} />
    </div>
  )
}

