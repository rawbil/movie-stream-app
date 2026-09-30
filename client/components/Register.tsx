"use client";

import { useFormik } from "formik";
import { Button } from "./ui/button";
import {
  Card,
  CardAction,
  CardContent,
  CardDescription,
  CardTitle,
} from "./ui/card";
import { Input } from "./ui/input";
import { Label } from "./ui/label";
import { RegisterSchema } from "@/lib/schema";
import { useMutation, useQuery } from "@tanstack/react-query";
import { RegisterUser } from "./mutations/auth.mutation";
import { useEffect, useState } from "react";
import {
  Eye,
  EyeClosed,
  EyeDashed,
  EyeOff,
  FileExclamationPointIcon,
} from "lucide-react";
import { GetMovieGenres } from "./queries/movie.queries";
import { ScrollArea } from "./ui/scroll-area";
import toast from "react-hot-toast";

type Values = {
  username: string;
  email: string;
  password: string;
  confirm_password: string;
};

type Genre = {
  genre_id: number;
  genre_name: string;
};

export default function Register() {
  const [errorMsg, setErrorMsg] = useState("");
  const [genreList, setGenreList] = useState<Genre[]>([]);
  const [selected_genres, setSelectedGenres] = useState<string[]>([]);
  const [isPassword, setIsPassword] = useState(true);
  const [iscPassword, setIsCPassword] = useState(true);

  const {
    error,
    data: genres,
    isFetching: genres_fetching,
  } = useQuery({
    queryKey: ["__get_genres"],
    queryFn: GetMovieGenres,
  });

  const registerMutation = useMutation({
    mutationFn: RegisterUser,
    onSuccess: (data) => {
      setErrorMsg("")
      toast.success(data.message)
    },
    onError: (error: any) => {
      setErrorMsg(
        error.message || (error.data as any).error || "error registering user",
      );
    },
  });

  useEffect(() => {
    if (genres) {
      setGenreList(genres.data.genres);
    }
  }, [genres]);

  useEffect(() => {
    console.log(selected_genres);
  }, [selected_genres]);

  const handleSelectedGenres = (
    e: React.ChangeEvent<HTMLInputElement>,
    genre: string,
  ) => {
    if (e.target.checked) {
      setSelectedGenres((prev_genres) => [...prev_genres, genre]);
    } else {
      setSelectedGenres((prev_genres) =>
        prev_genres.filter((item) => item != e.target.value),
      );
    }
  };

  const onSubmit = async (values: Values) => {
    if (selected_genres.length < 1) {
      setErrorMsg("choose at least one genre");
      return;
    }

    const { confirm_password, ...data } = values;
    const filtered_data = data;
    const payload = {
      ...filtered_data,
      fav_genres: selected_genres,
    };

    registerMutation.mutate(payload);
  };

  const { handleChange, handleBlur, handleSubmit, errors, values, touched } =
    useFormik({
      initialValues: {
        username: "",
        email: "",
        password: "",
        confirm_password: "",
      },
      validationSchema: RegisterSchema,
      onSubmit,
    });

  return (
    <section className="flex h-screen items-center justify-center">
      <Card className="w-150 max-w-full mx-auto">
        <CardTitle className="flex flex-col items-center">
          {errorMsg && (
            <p className="bg-red-500 p-2 rounded">
              <FileExclamationPointIcon />
              {errorMsg}
            </p>
          )}
          <h2 className="font-semibold text-xl">Create Account</h2>

          <div className="text-center">
            <span className="opacity-80">Already have an account?</span>{" "}
            <a href="/auth/login" className="underline hover:no-underline">
              Login
            </a>
          </div>
        </CardTitle>
        <CardContent>
          <form
            method="POST"
            className="flex flex-col gap-5"
            onSubmit={handleSubmit}
          >
            <div className="flex flex-col gap-2">
              <Label htmlFor="username">Username: </Label>
              <Input
                id="username"
                name="username"
                placeholder="enter your username"
                value={values.username}
                onChange={handleChange}
                onBlur={handleBlur}
              />
              {errors.username && touched.username && (
                <div className="text-red-500 text-sm">{errors.username}</div>
              )}
            </div>

            <div className="flex flex-col gap-2">
              <Label htmlFor="email">Email: </Label>
              <Input
                type="email"
                id="email"
                name="email"
                placeholder="youremail@gmail.com"
                value={values.email}
                onChange={handleChange}
                onBlur={handleBlur}
              />
              {errors.email && touched.email && (
                <div className="text-red-500 text-sm">{errors.email}</div>
              )}
            </div>

            <div className="flex flex-col gap-2">
              <Label htmlFor="password">Password: </Label>
              <div className="relative">
                <Input
                  type={isPassword ? "password" : "text"}
                  id="password"
                  name="password"
                  placeholder="enter your password"
                  value={values.password}
                  onChange={handleChange}
                  onBlur={handleBlur}
                />
                {isPassword ? (
                  <Eye
                    size={18}
                    className="absolute top-1/4 right-2 cursor-pointer"
                    onClick={() => setIsPassword(false)}
                  />
                ) : (
                  <EyeOff
                    size={18}
                    className="absolute top-1/4 right-2 cursor-pointer"
                    onClick={() => setIsPassword(true)}
                  />
                )}
              </div>
              {errors.password && (
                <div className="text-red-500 text-sm">{errors.password}</div>
              )}
            </div>

            <div className="flex flex-col gap-2">
              <Label htmlFor="confirm_password">Confirm Password: </Label>

              <div className="relative">
                <Input
                  type={iscPassword ? "password" : "text"}
                  id="confirm_password"
                  name="confirm_password"
                  placeholder="enter password again"
                  value={values.confirm_password}
                  onChange={handleChange}
                  onBlur={handleBlur}
                />
                {iscPassword ? (
                  <Eye
                    size={18}
                    className="absolute top-1/4 right-2 cursor-pointer"
                    onClick={() => setIsCPassword(false)}
                  />
                ) : (
                  <EyeOff
                    size={18}
                    className="absolute top-1/4 right-2 cursor-pointer"
                    onClick={() => setIsCPassword(true)}
                  />
                )}
              </div>
              {errors.confirm_password && (
                <div className="text-red-500 text-sm">
                  {errors.confirm_password}
                </div>
              )}
            </div>

            {/* genres */}
            <div>
              <hr />
              <h3 className="font-medium pt-4 pb-2 opacity-90">
                Choose Favourite Genres
              </h3>
              <ScrollArea className="h-25 rounded-md border">
                <ul>
                  {genreList ? (
                    genreList.length > 0 &&
                    genreList.map((genre: Genre) => (
                      <li key={genre.genre_id} className="p-1 flex gap-2">
                        <input
                          type="checkbox"
                          id={genre.genre_name}
                          value={genre.genre_name}
                          onChange={(e) =>
                            handleSelectedGenres(e, genre.genre_name)
                          }
                        />
                        <label
                          htmlFor={genre.genre_name}
                          className="capitalize"
                        >
                          {genre.genre_name}
                        </label>
                      </li>
                    ))
                  ) : (
                    <p>genre list emtpy</p>
                  )}
                </ul>
              </ScrollArea>
            </div>

            <Button
              type="submit"
              className="disabled:cursor-not-allowed"
              disabled={registerMutation.isPending}
            >
              {registerMutation.isPending ? "Submitting..." : "Submit"}
            </Button>
          </form>
        </CardContent>
      </Card>
    </section>
  );
}
