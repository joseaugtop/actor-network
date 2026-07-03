import { Component, OnInit, signal, computed, inject } from '@angular/core';
import { FormsModule } from '@angular/forms';
import { CommonModule } from '@angular/common';
import { CityService, CaminhoResponse } from './city.service';

@Component({
  selector: 'app-root',
  standalone: true,
  imports: [CommonModule, FormsModule],
  templateUrl: './app.html',
  styleUrl: './app.scss',
})
export class App implements OnInit {
  private cityService = inject(CityService);

  // Lista de capitais que alimenta os selects.
  readonly capitais = signal<string[]>([]);
  readonly loadingCapitais = signal<boolean>(true);
  readonly loadError = signal<string | null>(null);

  // Entradas do formulário.
  readonly origem = signal<string>('');
  readonly destino = signal<string>('');
  readonly combustivel = signal<number | null>(null); // preço do litro (R$)
  readonly autonomia = signal<number | null>(null); // km por litro

  // Estado da busca e resultado.
  readonly searching = signal<boolean>(false);
  readonly resultado = signal<CaminhoResponse | null>(null);
  readonly errorMessage = signal<string | null>(null);

  // Só habilita o botão quando todos os campos estão preenchidos.
  readonly canSearch = computed(
    () =>
      !this.searching() &&
      this.origem().trim().length > 0 &&
      this.destino().trim().length > 0 &&
      (this.combustivel() ?? 0) > 0 &&
      (this.autonomia() ?? 0) > 0,
  );

  ngOnInit(): void {
    this.cityService.getCapitais().subscribe({
      next: (res) => {
        this.capitais.set(res.capitais);
        this.loadingCapitais.set(false);
      },
      error: (err) => {
        this.loadingCapitais.set(false);
        this.loadError.set(
          'Não foi possível carregar a lista de capitais. Verifique se o backend está rodando em http://localhost:8081',
        );
        console.error(err);
      },
    });
  }

  // Inverte origem e destino.
  swap(): void {
    const a = this.origem();
    this.origem.set(this.destino());
    this.destino.set(a);
  }

  clearResults(): void {
    this.resultado.set(null);
    this.errorMessage.set(null);
  }

  // Dispara a busca do caminho mais barato no backend.
  buscar(): void {
    if (!this.canSearch()) return;
    this.clearResults();
    this.searching.set(true);
    this.cityService
      .caminho(
        this.origem().trim(),
        this.destino().trim(),
        this.combustivel()!,
        this.autonomia()!,
      )
      .subscribe({
        next: (res) => {
          this.resultado.set(res);
          this.searching.set(false);
        },
        error: (err) => {
          this.searching.set(false);
          this.errorMessage.set(
            err.status === 404
              ? 'Capital de origem ou destino não encontrada.'
              : 'Erro ao consultar o backend.',
          );
          console.error(err);
        },
      });
  }

  // Formata um número como moeda brasileira (R$).
  brl(value: number): string {
    return value.toLocaleString('pt-BR', {
      style: 'currency',
      currency: 'BRL',
    });
  }
}
