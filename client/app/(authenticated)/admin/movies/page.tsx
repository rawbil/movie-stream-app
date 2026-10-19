"use client";

import { useState, type FormEvent } from "react";
import { useMutation, useQuery, useQueryClient } from "@tanstack/react-query";
import toast from "react-hot-toast";
import { Pencil, Plus, Trash2, X } from "lucide-react";
import { CreateMovie, DeleteMovie, UpdateMovie } from "@/components/mutations/movie.mutation";
import { GetMovieGenres, GetMovies } from "@/components/queries/movie.queries";
import { Button } from "@/components/ui/button";
import { Input } from "@/components/ui/input";
import type { Movie } from "@/lib/types";

type MovieForm = {
  public_id?: string;
  imdb_id: string;
  title: string;
  poster_path: string;
  youtube_id: string;
  genres: string[];
};

const emptyForm: MovieForm = { imdb_id: "", title: "", poster_path: "", youtube_id: "", genres: [] };

export default function AdminMoviesPage() {
  const [page, setPage] = useState(1);
  const [form, setForm] = useState<MovieForm>(emptyForm);
  const queryClient = useQueryClient();
  const moviesQuery = useQuery({ queryKey: ["admin-movies", page], queryFn: () => GetMovies(page) });
  const genresQuery = useQuery({ queryKey: ["movie-genres"], queryFn: GetMovieGenres });
  const genres: { genre_name: string }[] = genresQuery.data?.data?.genres ?? [];
  const movies: Movie[] = moviesQuery.data?.data?.movies ?? [];
  const totalPages = Number(moviesQuery.data?.data?.total_pages ?? 1);

  const refreshList = async () => queryClient.invalidateQueries({ queryKey: ["admin-movies"] });
  const createMutation = useMutation({
    mutationFn: CreateMovie,
    onSuccess: async (data) => { toast.success(data.message); setForm(emptyForm); await refreshList(); },
    onError: (error: any) => toast.error(error.data?.error ?? "Could not create movie"),
  });
  const updateMutation = useMutation({
    mutationFn: UpdateMovie,
    onSuccess: async (data) => { toast.success(data.message); setForm(emptyForm); await refreshList(); },
    onError: (error: any) => toast.error(error.data?.error ?? "Could not update movie"),
  });
  const deleteMutation = useMutation({
    mutationFn: DeleteMovie,
    onSuccess: async (data) => { toast.success(data.message); await refreshList(); },
    onError: (error: any) => toast.error(error.data?.error ?? "Could not delete movie"),
  });

  const submit = (event: FormEvent<HTMLFormElement>) => {
    event.preventDefault();
    const payload = { ...form, youtube_id: form.youtube_id.trim() };
    if (form.public_id) updateMutation.mutate(payload as Required<MovieForm>);
    else createMutation.mutate(payload);
  };

  const editMovie = (movie: Movie) => setForm({
    public_id: movie.public_id,
    imdb_id: movie.imdb_id,
    title: movie.title,
    poster_path: movie.poster_path,
    youtube_id: movie.youtube_id.Valid ? movie.youtube_id.String : "",
    genres: movie.genres.Valid ? movie.genres.String.split(",").map((genre) => genre.trim()).filter(Boolean) : [],
  });

  const toggleGenre = (name: string) => setForm((current) => ({
    ...current,
    genres: current.genres.includes(name) ? current.genres.filter((genre) => genre !== name) : [...current.genres, name],
  }));

  return (
    <main className="mx-auto min-h-screen max-w-7xl px-4 py-8 sm:px-6">
      <header className="mb-8">
        <p className="text-sm font-medium uppercase tracking-wider text-muted-foreground">Administration</p>
        <h1 className="mt-2 text-3xl font-semibold">Manage movies</h1>
        <p className="mt-2 text-muted-foreground">Add titles, edit movie details and genres, or remove a title.</p>
      </header>

      <section className="mb-10 rounded-xl border border-border p-5">
        <div className="mb-5 flex items-center justify-between gap-4">
          <h2 className="text-xl font-semibold">{form.public_id ? "Edit movie" : "Add a movie"}</h2>
          {form.public_id && <Button type="button" variant="ghost" onClick={() => setForm(emptyForm)}><X className="mr-1 size-4" />Cancel edit</Button>}
        </div>
        <form onSubmit={submit} className="grid gap-4 sm:grid-cols-2">
          <label className="grid gap-1.5 text-sm">IMDb ID<Input required value={form.imdb_id} onChange={(e) => setForm({ ...form, imdb_id: e.target.value })} placeholder="tt0133093" /></label>
          <label className="grid gap-1.5 text-sm">Title<Input required minLength={3} maxLength={50} value={form.title} onChange={(e) => setForm({ ...form, title: e.target.value })} /></label>
          <label className="grid gap-1.5 text-sm sm:col-span-2">Poster URL<Input required type="url" value={form.poster_path} onChange={(e) => setForm({ ...form, poster_path: e.target.value })} placeholder="https://…" /></label>
          <label className="grid gap-1.5 text-sm sm:col-span-2">YouTube trailer ID<Input value={form.youtube_id} onChange={(e) => setForm({ ...form, youtube_id: e.target.value })} placeholder="Optional" /></label>
          <fieldset className="sm:col-span-2">
            <legend className="mb-2 text-sm font-medium">Genres</legend>
            {genresQuery.isLoading ? <p className="text-sm text-muted-foreground">Loading genres…</p> : (
              <div className="flex flex-wrap gap-2">
                {genres.map(({ genre_name }) => <label key={genre_name} className="flex items-center gap-2 rounded-lg border border-border px-3 py-2 text-sm"><input type="checkbox" checked={form.genres.includes(genre_name)} onChange={() => toggleGenre(genre_name)} />{genre_name}</label>)}
              </div>
            )}
            {genresQuery.isError && <p className="mt-2 text-sm text-destructive">Could not load genres. Refresh the page and try again.</p>}
            {!genresQuery.isLoading && !genresQuery.isError && genres.length === 0 && <p className="mt-2 text-sm text-muted-foreground">No genres are available yet.</p>}
          </fieldset>
          <div className="sm:col-span-2"><Button type="submit" disabled={!form.genres.length || createMutation.isPending || updateMutation.isPending}><Plus className="mr-1 size-4" />{form.public_id ? "Save changes" : "Create movie"}</Button></div>
        </form>
      </section>

      <section>
        <h2 className="mb-4 text-xl font-semibold">Movie catalog</h2>
        {moviesQuery.isLoading ? <p className="py-10 text-center text-muted-foreground">Loading movies…</p> : moviesQuery.isError ? <p className="py-10 text-center text-destructive">Could not load movies.</p> : (
          <div className="overflow-x-auto rounded-xl border border-border">
            <table className="w-full min-w-[700px] text-left text-sm">
              <thead className="bg-muted/50 text-muted-foreground"><tr><th className="p-3">Title</th><th className="p-3">IMDb</th><th className="p-3">Genres</th><th className="p-3 text-right">Actions</th></tr></thead>
              <tbody>{movies.map((movie) => <tr key={movie.public_id} className="border-t border-border"><td className="p-3 font-medium">{movie.title}</td><td className="p-3">{movie.imdb_id}</td><td className="p-3">{movie.genres.Valid ? movie.genres.String : "—"}</td><td className="p-3"><div className="flex justify-end gap-2"><Button type="button" size="sm" variant="outline" onClick={() => editMovie(movie)}><Pencil className="mr-1 size-3.5" />Edit</Button><Button type="button" size="sm" variant="destructive" disabled={deleteMutation.isPending} onClick={() => { if (window.confirm(`Delete “${movie.title}”?`)) deleteMutation.mutate(movie.public_id); }}><Trash2 className="mr-1 size-3.5" />Delete</Button></div></td></tr>)}</tbody>
            </table>
            {!movies.length && <p className="p-8 text-center text-muted-foreground">No movies on this page.</p>}
          </div>
        )}
        <div className="mt-4 flex items-center justify-between"><Button variant="outline" disabled={page <= 1} onClick={() => setPage(page - 1)}>Previous</Button><span className="text-sm text-muted-foreground">Page {page} of {totalPages}</span><Button variant="outline" disabled={page >= totalPages} onClick={() => setPage(page + 1)}>Next</Button></div>
      </section>
    </main>
  );
}
