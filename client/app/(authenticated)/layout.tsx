import AuthProtected from "@/components/Providers/Protected";

export default function AuthenticatedLayout({ children }: { children: React.ReactNode }) {
  return <AuthProtected>{children}</AuthProtected>;
}
