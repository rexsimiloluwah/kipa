export type SuccessResponse<T> = {
  status: boolean;
  message: string;
  data: T;
};

export type PageInfo = {
  total_items: number;
  total_pages: number;
  has_next_page: boolean;
  current_page: number;
};

export type PaginatedSuccessResponse<T> = SuccessResponse<T> & {
  page_info: PageInfo;
};
