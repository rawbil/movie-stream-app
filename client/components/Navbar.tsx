"use client";

import Link from "next/link";
import { Button } from "./ui/button";
import { useTheme } from "next-themes";
import {
  DropdownMenu,
  DropdownMenuContent,
  DropdownMenuItem,
  DropdownMenuTrigger,
} from "./ui/dropdown-menu";
import { Moon, Sun } from "lucide-react";
import { useMutation } from "@tanstack/react-query";
import { LogoutUser } from "./mutations/auth.mutation";
import { useAuthStore } from "@/app/auth/(store)/auth.store";
import { usePathname, useRouter } from "next/navigation";
import toast from "react-hot-toast";

export default function Navbar() {
  const { setTheme } = useTheme();
  const pathname = usePathname()

  const setAccessToken = useAuthStore((state) => state.setAccessToken);
  const router = useRouter();

  const logoutMutation = useMutation({
    mutationFn: LogoutUser,
    onSuccess: (data) => {
      setAccessToken("");
      router.push("/auth/login");
      toast.success(data.message);
    },
    onError: (error: any) => {
      toast.error(error.data.error);
      console.log(error.data.error);
    },
  });

  const logoutFn = async () => {
    logoutMutation.mutate();
  };

  return (
    <nav className="bg-sidebar-primary text-sidebar-primary-foreground flex justify-between items-center py-3 px-3 sticky top-0 z-50">
      <div>
        <h2>Movie Trailers</h2>
      </div>
      <ul className="flex items-center">
        <li className={`font-medium hover:underline px-3 ${pathname.endsWith("/") && "underline"}`}>
          <Link href={"/"}>Home</Link>
        </li>
        <li className={`font-medium hover:underline px-3 ${pathname.startsWith("/recommended") && "underline"}`}>
          <Link href={"/recommended"}>Recommended Movies</Link>
        </li>
      </ul>
      <div>
        <DropdownMenu>
          <DropdownMenuTrigger
            render={
              <Button variant="outline" size="icon">
                <Sun className="h-[1.2rem] w-[1.2rem] scale-100 rotate-0 transition-all dark:scale-0 dark:-rotate-90 text-foreground dark:text-background" />
                <Moon className="absolute h-[1.2rem] w-[1.2rem] scale-0 rotate-90 transition-all dark:scale-100 dark:rotate-0" />
                <span className="sr-only">Toggle theme</span>
              </Button>
            }
          />
          <DropdownMenuContent align="end">
            <DropdownMenuItem onClick={() => setTheme("light")}>
              Light
            </DropdownMenuItem>
            <DropdownMenuItem onClick={() => setTheme("dark")}>
              Dark
            </DropdownMenuItem>
            <DropdownMenuItem onClick={() => setTheme("system")}>
              System
            </DropdownMenuItem>
          </DropdownMenuContent>
        </DropdownMenu>
        <Button variant={"destructive"} onClick={logoutFn}>
          Logout
        </Button>
      </div>
    </nav>
  );
}
