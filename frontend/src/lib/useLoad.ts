import { useCallback, useEffect, useRef, useState } from "react";
import { errorMessage } from "../api/client";

// Loads data on mount (and when `key` changes) and exposes a reload function
// for after mutations.
export function useLoad<T>(load: () => Promise<T>, key: unknown = null) {
  const [data, setData] = useState<T | null>(null);
  const [error, setError] = useState<string | null>(null);
  const [version, setVersion] = useState(0);

  // Always call the latest `load` without re-running the effect on every render.
  const loadRef = useRef(load);
  useEffect(() => {
    loadRef.current = load;
  });

  useEffect(() => {
    let cancelled = false;
    loadRef.current().then(
      (result) => {
        if (!cancelled) {
          setData(result);
          setError(null);
        }
      },
      (err) => {
        if (!cancelled) setError(errorMessage(err));
      },
    );
    return () => {
      cancelled = true;
    };
  }, [key, version]);

  const reload = useCallback(() => setVersion((v) => v + 1), []);

  return { data, error, loading: data === null && error === null, reload };
}
