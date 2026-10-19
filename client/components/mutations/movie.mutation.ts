import { Axios } from "@/lib/axios";

type AddMovieReviewInput = {
  publicId: string;
  review: string;
};

export const AddMovieReview = async ({
  publicId,
  review,
}: AddMovieReviewInput) => {
  try {
    const response = await Axios.post(`/movies/add-review/${publicId}`, {
      review,
    });

    return response.data;
  } catch (error: any) {
    if (error.response) {
      throw error.response;
    }
    throw error;
  }
};

export async function CreateMovie(data: {
  imdb_id: string;
  title: string;
  poster_path: string;
  youtube_id: string;
  genres: string[];
}) {
  try {
    const response = await Axios.post("/movies/create", data);
    return response.data;
  } catch (error: any) {
    if (error.response) throw error.response;
    throw error;
  }
}

export async function UpdateMovie(data: {
  public_id: string;
  imdb_id: string;
  title: string;
  poster_path: string;
  youtube_id: string;
  genres: string[];
}) {
  try {
    const response = await Axios.put("/movies/update", data);
    return response.data;
  } catch (error: any) {
    if (error.response) throw error.response;
    throw error;
  }
}

export async function DeleteMovie(publicId: string) {
  try {
    const response = await Axios.delete(`/movies/delete/${publicId}`);
    return response.data;
  } catch (error: any) {
    if (error.response) throw error.response;
    throw error;
  }
}
