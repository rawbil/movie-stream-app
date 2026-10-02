"use client";

import type { Movie as MovieData } from "@/lib/types";
import { Card } from "./ui/card";
import Image from "next/image";
import { Play } from "lucide-react";
import { useAuthStore } from "@/app/auth/(store)/auth.store";
import { Button } from "./ui/button";

interface MovieProps {
  movie: MovieData;
}

export default function Movie({ movie }: MovieProps) {
  const accessToken = useAuthStore((state) => state.accessToken);

  const isAuthenticated = !!accessToken;

  const genres = movie.genres.Valid
    ? movie.genres.String.split(",")
        .map((genre) => genre.trim())
        .filter(Boolean)
    : [];
  const hasTrailer = movie.youtube_id.Valid && movie.youtube_id.String.trim();

  return (
    <li className="min-w-0 list-none">
      <Card className="h-full gap-0 overflow-hidden rounded-xl border border-border/70 bg-card py-0 shadow-sm transition duration-200 hover:-translate-y-1 hover:shadow-lg">
        <div className="relative aspect-[2/3] w-full overflow-hidden bg-muted">
          <Image
            src={movie.poster_path}
            alt={movie.title}
            fill
            sizes="(max-width: 640px) 50vw, (max-width: 1024px) 33vw, (max-width: 1280px) 25vw, 20vw"
            className="object-cover transition duration-300 group-hover/card:scale-105"
          />
          {movie.ranking_name.Valid && movie.ranking_name.String && (
            <span className="absolute left-3 top-3 max-w-[calc(100%-1.5rem)] truncate rounded-md bg-background/90 px-2 py-1 text-xs font-semibold text-foreground shadow-sm backdrop-blur-sm">
              {movie.ranking_name.String}
            </span>
          )}
        </div>

        <div className="flex min-h-48 flex-col gap-3 p-4">
          <div className="min-w-0">
            <h3 className="line-clamp-2 text-base font-semibold leading-tight tracking-normal">
              {movie.title}
            </h3>
            <p className="mt-1 truncate text-xs text-muted-foreground">
              IMDb {movie.imdb_id}
            </p>
          </div>

          {genres.length > 0 && (
            <ul className="flex flex-wrap gap-1.5" aria-label="Genres">
              {genres.map((genre) => (
                <li
                  className="rounded-md bg-secondary px-2 py-1 text-[11px] font-medium text-secondary-foreground"
                  key={genre}
                >
                  {genre}
                </li>
              ))}
            </ul>
          )}

          <div className="mt-auto flex items-center justify-between gap-3 border-t border-border/60 pt-3">
            <span className="truncate text-xs text-muted-foreground">
              <Button>Add Review</Button>
            </span>

            {hasTrailer && (
              <a
                className="inline-flex shrink-0 items-center gap-1.5 text-xs font-semibold text-primary underline-offset-4 transition hover:underline"
                href={
                  isAuthenticated
                    ? `https://www.youtube.com/watch?v=${movie.youtube_id.String}`
                    : "/auth/login"
                }
                target={isAuthenticated ? "_blank" : undefined}
                rel={isAuthenticated ? "noopener noreferrer" : undefined}
                aria-label={`Watch trailer for ${movie.title}`}
              >
                <Play aria-hidden="true" className="size-3.5 fill-current" />
                Trailer
              </a>
            )}
          </div>
        </div>
      </Card>
    </li>
  );
}
