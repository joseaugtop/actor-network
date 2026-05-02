import { Injectable, inject } from '@angular/core';
import { HttpClient, HttpParams } from '@angular/common/http';
import { Observable } from 'rxjs';

export interface ActorsResponse {
  count: number;
  actors: string[];
}

export interface BfsResponse {
  path: string[];
  length: number;
  found: boolean;
  message?: string;
}

export interface PathItem {
  path: string[];
  length: number;
}

export interface Bfs8Response {
  paths: PathItem[];
  count: number;
  minLength?: number;
  maxLength?: number;
  truncated?: boolean;
  cap?: number;
  found: boolean;
  message?: string;
}

export interface ShowResponse {
  count: number;
  adjacency: Record<string, string[]>;
}

@Injectable({ providedIn: 'root' })
export class ActorService {
  private http = inject(HttpClient);

  getActors(): Observable<ActorsResponse> {
    return this.http.get<ActorsResponse>('/actors');
  }

  bfs(from: string, to: string): Observable<BfsResponse> {
    const params = new HttpParams().set('from', from).set('to', to);
    return this.http.get<BfsResponse>('/bfs', { params });
  }

  bfs8(from: string, to: string, max = 8): Observable<Bfs8Response> {
    const params = new HttpParams()
      .set('from', from)
      .set('to', to)
      .set('max', String(max));
    return this.http.get<Bfs8Response>('/bfs8', { params });
  }

  show(): Observable<ShowResponse> {
    return this.http.get<ShowResponse>('/show');
  }
}
