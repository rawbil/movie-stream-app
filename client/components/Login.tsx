"use client";

import { useFormik } from "formik";
import { Card, CardContent, CardTitle } from "./ui/card";
import { Input } from "./ui/input";
import { Label } from "./ui/label";
import { useState } from "react";
import { Eye, EyeOff, FileExclamationPointIcon } from "lucide-react";
import { LoginSchema } from "@/lib/schema";
import { Button } from "./ui/button";
import { useMutation } from "@tanstack/react-query";
import { LoginUser } from "./mutations/auth.mutation";
import { useAuthStore } from "@/app/auth/(store)/auth.store";
import toast from "react-hot-toast";
import { useRouter } from "next/navigation";
import { useUserStore } from "@/app/auth/(store)/user.store";

type Values = {
  email: string;
  password: string;
};

export default function Login() {
  const [isPassword, setIsPassword] = useState(true);
  const [errorMsg, setErrorMsg] = useState("");
  const setAccessToken = useAuthStore((state) => state.setAccessToken);
  const setUser = useUserStore((state) => state.setUser);

  const router = useRouter();

  const loginMutation = useMutation({
    mutationFn: LoginUser,
    onSuccess: (data) => {
      toast.success(data.message);
      setErrorMsg("");

      setAccessToken(data.access_token);
      setUser({
        public_id: data.user.public_id,
        email: data.user.email,
        created_at: data.user.created_at,
        updated_at: data.user.updated_at,
      });

      router.push("/");
    },
    onError: (error: any) => {
      setErrorMsg(error.data.error || error.message || "error logging in");
    },
  });

  const onSubmit = async (values: Values) => {
    loginMutation.mutate(values);
  };

  const { handleChange, handleBlur, handleSubmit, errors, values, touched } =
    useFormik({
      initialValues: {
        email: "",
        password: "",
      },
      validationSchema: LoginSchema,
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
          <h2 className="font-semibold text-xl">Welcome Back, Login</h2>

          <div className="text-center">
            <span className="opacity-80">Don't have an account?</span>{" "}
            <a href="/auth/register" className="underline hover:no-underline">
              Register
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

            <Button
              type="submit"
              className="disabled:cursor-not-allowed"
              disabled={loginMutation.isPending}
            >
              {loginMutation.isPending ? "Submitting..." : "Submit"}
            </Button>
          </form>
        </CardContent>
      </Card>
    </section>
  );
}
