import AdminProtected from "@/components/Providers/AdminProtected";

export default function AdminLayout({ children }: { children: React.ReactNode }) {
  return <AdminProtected>{children}</AdminProtected>;
}
