import Reviews from "@/components/Reviews";

export default async function ReviewsPage({
  params,
}: {
  params: Promise<{ public_id: string }>;
}) {
  const { public_id } = await params;

  return <Reviews publicId={public_id} />;
}
