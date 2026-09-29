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
