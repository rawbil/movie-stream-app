import * as yup from "yup";
const passwordRegex = /^(?=.{6,20}$)(?=.*[A-Za-z])(?=.*\d)(?=.*[^A-Za-z\d]).*$/;
const emailRegex =
  /^(([^<>()[\]\\.,;:\s@"]+(\.[^<>()[\]\\.,;:\s@"]+)*)|(".+"))@((\[[0-9]{1,3}\.[0-9]{1,3}\.[0-9]{1,3}\.[0-9]{1,3}\])|(([a-zA-Z\-0-9]+\.)+[a-zA-Z]{2,}))$/;

export type RegisterParams = {
  username: string;
  email: string;
  password: string;
  fav_genres: string[];
};

export type LoginParams = {
  email: string;
  password: string;
};

export const RegisterSchema = yup.object().shape({
  username: yup
    .string()
    .required("username field required")
    .min(3, "username should be at least 3 characters")
    .max(20, "username should not be more than 20 characters")
    .trim(),

  email: yup
    .string()
    .required("email field required")
    .matches(emailRegex, "invalid email format")
    .trim()
    .lowercase(),

  password: yup
    .string()
    .required("password field required")
    .trim()
    .matches(
      passwordRegex,
      "password should be at between 6-20 characters long, have at least alphanumerical and have at least one special character",
    ),
  confirm_password: yup
    .string()
    .required("confirm_password field is required")
    .oneOf([yup.ref("password")], "passwords do not match"),
});

export const LoginSchema = yup.object().shape({
  email: yup
    .string()
    .required("email field required")
    .matches(emailRegex, "invalid email format")
    .trim()
    .lowercase(),

  password: yup
    .string()
    .required("password field required")
    .trim()
    .matches(
      passwordRegex,
      "password should be at between 6-20 characters long, have at least alphanumerical and have at least one special character",
    ),
});
