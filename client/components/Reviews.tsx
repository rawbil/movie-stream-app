"use client";

import AuthProvider from "./Providers/AuthProvider";
import {
  Dialog,
  DialogContent,
  DialogDescription,
  DialogHeader,
  DialogTitle,
  DialogTrigger,
} from "@/components/ui/dialog";
import { Button } from "./ui/button";
import Image from "next/image";
import { useSearchParams } from "next/navigation";
import { useMutation, useQuery, useQueryClient } from "@tanstack/react-query";
import { AddMovieReview } from "./mutations/movie.mutation";
import toast from "react-hot-toast";
import { useState } from "react";
import { GetMovieReviews } from "./queries/movie.queries";
import {
  ChevronLeft,
  ChevronRight,
  MessageSquareText,
  PenLine,
} from "lucide-react";

type MovieReview = {
  review: string;
  created_at: string;
  username: string;
  movie_public_id: string;
};

type MovieReviewsResponse = {
  data: {
    reviews: MovieReview[];
    page: number;
    total_pages: number;
    total_reviews: number;
  };
};

export default function Reviews({ publicId }: { publicId: string }) {
  const searchParams = useSearchParams();
  const posterPath = searchParams.get("poster_path") ?? "";
  const title = searchParams.get("title") ?? "Movie review";
  const [review, setReview] = useState("");
  const [isDialogOpen, setIsDialogOpen] = useState(false);
  const [page, setPage] = useState(1);
  const queryClient = useQueryClient();

  const addReviewMutation = useMutation({
    mutationFn: AddMovieReview,
    onSuccess: () => {
      setReview("");
      setIsDialogOpen(false);
      setPage(1);
      queryClient.invalidateQueries({ queryKey: ["reviews", publicId] });
      toast.success("Review added");
    },
    onError: () => {
      toast.error("Unable to add your review");
    },
  });

  const {
    data: reviewsData,
    error: reviewsError,
    isLoading: isLoadingReviews,
    isFetching: isFetchingReviews,
  } = useQuery<MovieReviewsResponse>({
    queryFn: () => GetMovieReviews(publicId, page),
    queryKey: ["reviews", publicId, page],
  });

  const reviews = reviewsData?.data.reviews ?? [];
  const totalReviews = reviewsData?.data.total_reviews ?? 0;
  const totalPages = Math.max(1, reviewsData?.data.total_pages ?? 1);

  function handleSubmit(e: React.SubmitEvent) {
    e.preventDefault();
    addReviewMutation.mutate({ publicId, review: review.trim() });
  }

  return (
    <AuthProvider>
      <main className="min-h-screen bg-background">
        <section className="border-b border-border bg-muted/30">
          <div className="mx-auto grid max-w-5xl gap-6 px-4 py-8 sm:grid-cols-[10rem_minmax(0,1fr)] sm:items-end sm:px-6 lg:px-8">
            <div className="relative mx-auto aspect-[2/3] w-32 overflow-hidden rounded-lg border border-border bg-muted shadow-sm sm:mx-0 sm:w-40">
              {posterPath ? (
                <Image
                  src={posterPath}
                  alt={`${title} poster`}
                  fill
                  sizes="(max-width: 640px) 8rem, 10rem"
                  className="object-cover"
                />
              ) : (
                <div className="flex h-full items-center justify-center p-4 text-center text-xs text-muted-foreground">
                  No poster available
                </div>
              )}
            </div>

            <div className="min-w-0">
              <p className="text-xs font-semibold uppercase tracking-[0.12em] text-muted-foreground">
                Movie discussion
              </p>
              <h1 className="mt-2 text-3xl font-semibold tracking-normal sm:text-4xl">
                {title}
              </h1>
              <p className="mt-3 max-w-2xl text-sm leading-6 text-muted-foreground">
                Read what others thought, then add your own take.
              </p>

              <Dialog open={isDialogOpen} onOpenChange={setIsDialogOpen}>
                <DialogTrigger
                  render={<Button className="mt-5" size="lg" />}
                >
                  <PenLine aria-hidden="true" />
                  Write a review
                </DialogTrigger>
                <DialogContent className="rounded-xl">
                  <DialogHeader>
                    <DialogTitle>Share your thoughts</DialogTitle>
                    <DialogDescription>
                      Your review helps other viewers decide what to watch.
                    </DialogDescription>
                  </DialogHeader>
                  <form className="grid gap-4" onSubmit={handleSubmit}>
                    <label className="grid gap-2 text-sm font-medium" htmlFor="review">
                      Your review
                      <textarea
                        id="review"
                        name="review"
                        placeholder="What stood out to you?"
                        className="min-h-32 w-full resize-y rounded-lg border border-input bg-background px-3 py-2 text-sm outline-none transition placeholder:text-muted-foreground focus-visible:border-ring focus-visible:ring-3 focus-visible:ring-ring/30"
                        required
                        value={review}
                        onChange={(e) => setReview(e.target.value)}
                      />
                    </label>
                    <Button
                      type="submit"
                      className="justify-self-end"
                      disabled={addReviewMutation.isPending || !review.trim()}
                    >
                      {addReviewMutation.isPending ? "Submitting..." : "Publish review"}
                    </Button>
                  </form>
                </DialogContent>
              </Dialog>
            </div>
          </div>
        </section>

        <section className="mx-auto w-full max-w-3xl px-4 py-10 sm:px-6 lg:px-8">
        <div className="mb-6 flex items-end justify-between gap-4 border-b border-border pb-4">
          <div>
            <div className="flex items-center gap-2 text-xs font-semibold uppercase tracking-[0.12em] text-muted-foreground">
              <MessageSquareText aria-hidden="true" className="size-3.5" />
              Community
            </div>
            <h2 className="mt-2 text-2xl font-semibold tracking-normal">Reviews</h2>
          </div>
          {!reviewsError && (
            <p className="shrink-0 text-sm text-muted-foreground">
              {totalReviews} {totalReviews === 1 ? "review" : "reviews"}
            </p>
          )}
        </div>

        {isLoadingReviews ? (
          <p className="py-10 text-center text-sm text-muted-foreground">
            Loading reviews...
          </p>
        ) : reviewsError ? (
          <div className="border border-destructive/30 bg-destructive/10 px-4 py-3 text-sm text-destructive">
            Unable to load reviews right now.
          </div>
        ) : reviews.length > 0 ? (
          <ul className="divide-y divide-border border-y border-border" aria-busy={isFetchingReviews}>
            {reviews.map((movieReview, index) => (
              <li
                key={`${movieReview.username}-${movieReview.created_at}-${index}`}
                className="py-6"
              >
                <div className="flex items-center justify-between gap-3">
                  <div className="flex min-w-0 items-center gap-2.5">
                    <span
                      className="flex size-8 shrink-0 items-center justify-center rounded-full bg-secondary text-xs font-semibold text-secondary-foreground"
                      aria-hidden="true"
                    >
                      {movieReview.username.slice(0, 1).toUpperCase()}
                    </span>
                    <p className="truncate font-medium">{movieReview.username}</p>
                  </div>
                  <time
                    className="shrink-0 text-xs text-muted-foreground"
                    dateTime={movieReview.created_at}
                  >
                    {new Intl.DateTimeFormat(undefined, {
                      dateStyle: "medium",
                    }).format(new Date(movieReview.created_at))}
                  </time>
                </div>
                <p className="mt-3 whitespace-pre-wrap text-sm leading-6 text-muted-foreground sm:pl-10.5">
                  {movieReview.review}
                </p>
              </li>
            ))}
          </ul>
        ) : (
          <div className="border border-dashed border-border px-6 py-12 text-center">
            <p className="font-medium">No reviews yet</p>
            <p className="mt-1 text-sm text-muted-foreground">
              Be the first to share your thoughts on this movie.
            </p>
          </div>
        )}

        {totalPages > 1 && (
          <nav
            className="mt-6 flex items-center justify-between gap-4"
            aria-label="Reviews pagination"
          >
            <Button
              variant="outline"
              size="icon"
              onClick={() => setPage((currentPage) => currentPage - 1)}
              disabled={page === 1 || isFetchingReviews}
              aria-label="Previous reviews page"
            >
              <ChevronLeft aria-hidden="true" />
            </Button>
            <span className="text-sm text-muted-foreground">
              Page {page} of {totalPages}
            </span>
            <Button
              variant="outline"
              size="icon"
              onClick={() => setPage((currentPage) => currentPage + 1)}
              disabled={page === totalPages || isFetchingReviews}
              aria-label="Next reviews page"
            >
              <ChevronRight aria-hidden="true" />
            </Button>
          </nav>
        )}
        </section>
      </main>
    </AuthProvider>
  );
}
