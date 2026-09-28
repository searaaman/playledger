const rupees = new Intl.NumberFormat("en-IN", {
  style: "currency",
  currency: "INR",
  maximumFractionDigits: 2,
  minimumFractionDigits: 0,
});

export function money(amount: number): string {
  // Avoid showing "-₹0" for tiny float leftovers like -0.0000001.
  const rounded = Math.round(amount * 100) / 100;
  return rupees.format(Object.is(rounded, -0) ? 0 : rounded);
}

// Treat anything under half a paisa as zero when deciding who owes.
export function isZero(amount: number): boolean {
  return Math.abs(amount) < 0.005;
}

export function time(iso: string): string {
  return new Date(iso).toLocaleTimeString("en-IN", {
    hour: "numeric",
    minute: "2-digit",
  });
}

export function date(iso: string): string {
  return new Date(iso).toLocaleDateString("en-IN", {
    weekday: "short",
    day: "numeric",
    month: "short",
    year: "numeric",
  });
}

export function shortDate(iso: string): string {
  return new Date(iso).toLocaleDateString("en-IN", {
    day: "numeric",
    month: "short",
  });
}

export function timeRange(startIso: string, endIso: string): string {
  return `${time(startIso)} – ${time(endIso)}`;
}

export function dateTime(iso: string): string {
  return `${shortDate(iso)}, ${time(iso)}`;
}

// Value for <input type="date"> in the browser's local time zone.
export function todayInputValue(): string {
  const now = new Date();
  const offset = now.getTimezoneOffset() * 60_000;
  return new Date(now.getTime() - offset).toISOString().slice(0, 10);
}

// Combines a local date ("2026-09-27") and time ("18:00") into an ISO timestamp.
export function toIso(dateValue: string, timeValue: string): string {
  return new Date(`${dateValue}T${timeValue}`).toISOString();
}
