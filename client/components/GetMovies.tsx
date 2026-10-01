"use client";

import { useQuery } from "@tanstack/react-query";
import { GetMovies } from "./queries/movie.queries";
import ServerLoading from "./Loaders/ServerLoading";
import { Movie as MovieType } from "@/lib/types";
import Movie from "./Movie";
import { useRouter, useSearchParams } from "next/navigation";
import { useEffect, useState } from "react";
import { Button } from "./ui/button";

export default function GetAllMovies() {
  const [page, setPage] = useState(1);
  const searchParams = useSearchParams();
  const router = useRouter();

  const pageParam = Number(searchParams.get("page") || 1);

  const {
    data: moviesData,
    error: moviesError,
    isFetching: moviesFetching,
  } = useQuery({
    queryFn: () => GetMovies(page),
    queryKey: ["all-movies", page],
  });

  useEffect(() => {
    setPage(pageParam);
  }, [pageParam]);

  if (moviesFetching) return <ServerLoading />;

  if (moviesError) {
    return (
      <main className="mx-auto flex min-h-[50vh] max-w-7xl items-center justify-center px-4 py-12">
        <p className="rounded-lg border border-destructive/30 bg-destructive/10 px-4 py-3 text-sm text-destructive">
          Unable to load movies right now.
        </p>
      </main>
    );
  }

  const movies: MovieType[] = moviesData?.data?.movies ?? [];
  const totalPages = moviesData?.data?.total_pages;

  const handlePageChange = (newPage: number) => {
    setPage(newPage);
    const params = new URLSearchParams()
    params.set("page", String(newPage))
    router.push(`?${params.toString()}`);
  };

  return (
    <main className="min-h-screen bg-background px-4 py-8 sm:px-6 lg:px-8">
      <div className="mx-auto max-w-7xl">
        <header className="mb-6 flex items-end justify-between gap-4">
          <div>
            <p className="text-xs font-semibold uppercase tracking-[0.12em] text-muted-foreground">
              Library
            </p>
            <h2 className="mt-1 text-2xl font-semibold tracking-normal sm:text-3xl">
              Movies
            </h2>
          </div>
          <p className="text-sm text-muted-foreground">
            {movies.length} {movies.length === 1 ? "title" : "titles"}
          </p>
        </header>

        {movies.length > 0 ? (
          <ul className="grid list-none grid-cols-2 gap-4 p-0 sm:grid-cols-3 sm:gap-5 lg:grid-cols-4 xl:grid-cols-5">
            {movies.map((movie) => (
              <Movie key={movie.public_id} movie={movie} />
            ))}
          </ul>
        ) : (
          <div className="rounded-xl border border-dashed border-border px-6 py-16 text-center">
            <p className="font-medium">No movies found</p>
            <p className="mt-1 text-sm text-muted-foreground">
              Your movie library is empty.
            </p>
          </div>
        )}
      </div>

      {/* Pagination */}
      <div className="flex justify-between items-center mt-10 border-t pt-5">
        <Button
          onClick={() => handlePageChange(Math.max(1, page - 1))}
          disabled={page === 1}
          className="disabled:bg-transparent disabled:text-white  disabled:cursor-not-allowed"
        >
          Prev
        </Button>
        <span className="opacity-70">
          Page {page} of {totalPages}
        </span>
        <Button
          onClick={() => handlePageChange(Math.min(totalPages, page + 1))}
          disabled={page === totalPages}
          className="disabled:bg-transparent disabled:text-white disabled:cursor-not-allowed"
        >
          Next
        </Button>
      </div>
    </main>
  );
}
