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
