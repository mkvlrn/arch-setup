#!/usr/bin/env bash
# STEP_NAME is read by install.sh after this file is sourced.
# shellcheck disable=SC2034

STEP_NAME='Configuring proton drive mount with rclone'

step_run() {
  [[ "${CI:-}" == "true" ]] && return 0

  if ! rclone listremotes | grep -qx 'proton:'; then
    local proton_secrets="$HOME/.config/rclone/proton-secrets"

    if [[ -f "$proton_secrets" ]]; then
      # shellcheck disable=SC1090
      . "$proton_secrets"

      if [[ -z "${PROTON_USER:-}" ||
        -z "${PROTON_PASSWORD:-}" ||
        -z "${PROTON_TOTP:-}" ]]; then
        printf 'Missing Proton credentials in %s\n' "$proton_secrets" >&2
        return 1
      fi

      local otp_secret
      otp_secret=$(rclone obscure "$PROTON_TOTP")

      run rclone config create proton protondrive \
        username="$PROTON_USER" \
        password="$PROTON_PASSWORD" \
        otp_secret_key="$otp_secret"

      unset PROTON_USER PROTON_PASSWORD PROTON_TOTP otp_secret
    fi
  fi
}
