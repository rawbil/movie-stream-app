"use client";

import { useQuery } from "@tanstack/react-query";
import AuthProtected from "./Providers/Protected";
import { GetRecommendedMovies } from "./queries/movie.queries";
import ServerLoading from "./Loaders/ServerLoading";
import Movie from "./Movie";
import type { Movie as MovieType } from "@/lib/types";

export default function RecommendedMovies() {
  const {
    data: moviesData,
    error,
    isFetching,
  } = useQuery({
    queryFn: GetRecommendedMovies,
    queryKey: ["recommended"],
  });

  if (isFetching) return <ServerLoading />;

  const movies: MovieType[] = moviesData?.movies ?? [];

  return (
    <AuthProtected>
      <main className="min-h-screen bg-background px-4 py-8 sm:px-6 lg:px-8">
        <div className="mx-auto max-w-7xl">
          <header className="mb-8 flex flex-col gap-3 border-b border-border pb-6 sm:flex-row sm:items-end sm:justify-between">
            <div>
              <p className="text-xs font-semibold uppercase tracking-[0.12em] text-muted-foreground">
                Curated for you
              </p>
              <h1 className="mt-1 text-2xl font-semibold tracking-normal sm:text-3xl">
                Recommended movies
              </h1>
              <p className="mt-2 max-w-2xl text-sm text-muted-foreground">
                A selection of films picked from your library.
              </p>
            </div>
            {!error && (
              <p className="text-sm text-muted-foreground">
                {movies.length} {movies.length === 1 ? "title" : "titles"}
              </p>
            )}
          </header>

          {error ? (
            <div className="border border-destructive/30 bg-destructive/10 px-4 py-3 text-sm text-destructive">
              Unable to load recommendations right now.
            </div>
          ) : movies.length > 0 ? (
            <ul className="grid list-none grid-cols-2 gap-4 p-0 sm:grid-cols-3 sm:gap-5 lg:grid-cols-4 xl:grid-cols-5">
              {movies.map((movie) => (
                <Movie key={movie.public_id} movie={movie} />
              ))}
            </ul>
          ) : (
            <div className="border border-dashed border-border px-6 py-16 text-center">
              <p className="font-medium">No recommendations yet</p>
              <p className="mt-1 text-sm text-muted-foreground">
                Your movie genres have not been reviewed yet.
              </p>
            </div>
          )}
        </div>
      </main>
    </AuthProtected>
  );
}
