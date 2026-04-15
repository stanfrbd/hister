export interface HistoryItem {
  id?: number;
  query: string;
  url: string;
  title: string;
  updated_at?: string;
  added?: number;
  favicon?: string;
  text?: string;
}
