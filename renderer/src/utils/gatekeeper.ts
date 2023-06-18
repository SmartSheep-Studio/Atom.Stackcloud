import { type RouteLocationNormalizedLoaded } from "vue-router"
import { useAccount } from "@/stores/account"
import { useSnackbar } from "@/stores/snackbar"
import wildcard from "wildcard-match"

function notify() {
  const $snackbar = useSnackbar()

  $snackbar.show({
    text: "You are accessing a page that requires authorization. Please sign in or try again with an account with a higher privilege level.",
    color: "warning"
  })
}

export function hasUserPermissions(...requires: string[]) {
  const $account = useAccount()
  if (!$account.isLoggedIn || $account.profile == null) {
    return false
  }

  for (const require of requires) {
    let passed = false
    for (const perm of $account.profile.permissions ?? []) {
      if (wildcard(perm)(require)) {
        passed = true
        break
      }
    }

    if (!passed) {
      return false
    }
  }

  return true
}

export function keepGate(to: RouteLocationNormalizedLoaded, notification = true) {
  const $account = useAccount()
  const meta: any = to?.meta?.gatekeeper ?? {}

  if (meta?.must === true && !$account.isLoggedIn) {
    if (notification) {
      notify()
    }
    return false
  } else if (meta?.permissions != null && !hasUserPermissions(...meta?.permissions)) {
    if (notification) {
      notify()
    }
    return false
  }

  return true
}
