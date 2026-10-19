import Link from "next/link";
import { Button } from "@/components/ui/button";

export default function AdminPage() {
  return (
    <main className="mx-auto min-h-[70vh] max-w-6xl px-4 py-10">
      <p className="text-sm font-medium uppercase tracking-wider text-muted-foreground">Administration</p>
      <h1 className="mt-2 text-3xl font-semibold">Admin dashboard</h1>
      <p className="mt-2 text-muted-foreground">Manage the movie catalog.</p>
      <Button className="mt-6" render={<Link href="/admin/movies" />}>Manage movies</Button>
    </main>
  );
}
