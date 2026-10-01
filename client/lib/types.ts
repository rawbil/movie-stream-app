export interface NullString {
  String: string;
  Valid: boolean;
}

export interface NullInt32 {
  Int32: number;
  Valid: boolean;
}

export interface Movie {
  public_id: string;
  imdb_id: string;
  title: string;
  poster_path: string;
  youtube_id: NullString;
  admin_review: NullString;
  ranking_value: NullInt32;
  ranking_name: NullString;
  genres: NullString;
}
