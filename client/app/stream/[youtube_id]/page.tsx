"use client";

import dynamic from "next/dynamic";
import { useParams, useRouter } from "next/navigation";
import { useState } from "react";
import { ArrowLeft, CircleAlert, ExternalLink } from "lucide-react";
import { Button } from "@/components/ui/button";

const ReactPlayer = dynamic(() => import("react-player"), {
  ssr: false,
  loading: () => <div className="h-full w-full animate-pulse bg-muted" />,
});

export default function StreamMoviePage() {
  const router = useRouter();
  const params = useParams<{ youtube_id: string }>();
  const [hasPlaybackError, setHasPlaybackError] = useState(false);
  const youtubeId = params.youtube_id;
  const videoUrl = `https://www.youtube.com/watch?v=${encodeURIComponent(youtubeId)}`;

  return (
    <main className="min-h-screen bg-background px-4 py-6 sm:px-6 lg:px-8">
      <div className="mx-auto max-w-6xl">
        <header className="mb-6 flex items-center justify-between gap-4">
          <Button variant="ghost" onClick={() => router.back()}>
            <ArrowLeft aria-hidden="true" />
            Back
          </Button>
          <a
            className="inline-flex items-center gap-1.5 text-sm font-medium text-muted-foreground transition hover:text-foreground"
            href={videoUrl}
            target="_blank"
            rel="noopener noreferrer"
          >
            Open on YouTube
            <ExternalLink aria-hidden="true" className="size-3.5" />
          </a>
        </header>

        <section aria-label="Trailer player">
          <div className="relative aspect-video w-full overflow-hidden rounded-lg border border-border bg-black shadow-sm">
            {hasPlaybackError ? (
              <div className="flex h-full flex-col items-center justify-center gap-3 px-6 text-center text-white">
                <CircleAlert aria-hidden="true" className="size-7" />
                <p className="text-sm">This trailer is unavailable for embedded playback.</p>
                <a
                  className="inline-flex items-center gap-1.5 text-sm font-medium underline underline-offset-4"
                  href={videoUrl}
                  target="_blank"
                  rel="noopener noreferrer"
                >
                  Watch on YouTube
                  <ExternalLink aria-hidden="true" className="size-3.5" />
                </a>
              </div>
            ) : (
              <ReactPlayer
                src={videoUrl}
                controls
                playsInline
                playing={true}
                width="100%"
                height="100%"
                onError={() => setHasPlaybackError(true)}
              />
            )}
          </div>
        </section>
      </div>
    </main>
  );
}
