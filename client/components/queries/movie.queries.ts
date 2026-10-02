import { Axios } from "@/lib/axios";

export async function GetMovieGenres() {
  try {
    const response = await Axios.get("/movies/genres/list");
    return response.data;
  } catch (error: any) {
    if (error.response) {
      throw error.response;
    }
    throw error;
  }
}

export async function GetMovies(page: number) {
  try {
    const response = await Axios.get(`/movies/all?page=${page}`);
    return response.data;
  } catch (error: any) {
    if (error.response) {
      throw error.response;
    }
    throw error;
  }
}

export async function GetRecommendedMovies() {
  try {
    const response = await Axios.get("/movies/recommended");
    return response.data;
  } catch (error: any) {
    if (error.response) {
      throw error.response;
    }
    throw error;
  }
}
