import { Injectable, inject } from '@angular/core';
import { HttpClient, HttpParams } from '@angular/common/http';
import { Observable } from 'rxjs';

// Resposta de GET /capitais — lista usada nos selects.
export interface CapitaisResponse {
  count: number;
  capitais: string[];
}

// Resposta de GET /caminho — espelha o struct Result do backend.
export interface CaminhoResponse {
  path: string[]; // capitais da origem ao destino
  distance: number; // distância total em km
  fuelCost: number; // gasto com combustível
  tollCost: number; // gasto com pedágios
  totalCost: number; // combustível + pedágios
  found: boolean; // existe rota?
  message?: string; // mensagem quando não há rota
}

@Injectable({ providedIn: 'root' })
export class CityService {
  private http = inject(HttpClient);

  getCapitais(): Observable<CapitaisResponse> {
    return this.http.get<CapitaisResponse>('/capitais');
  }

  caminho(
    origem: string,
    destino: string,
    combustivel: number,
    autonomia: number,
  ): Observable<CaminhoResponse> {
    const params = new HttpParams()
      .set('origem', origem)
      .set('destino', destino)
      .set('combustivel', String(combustivel))
      .set('autonomia', String(autonomia));
    return this.http.get<CaminhoResponse>('/caminho', { params });
  }
}
